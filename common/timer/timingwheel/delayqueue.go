package timingwheel

import (
	"context"
	"github.com/ahKevinXy/go-web-tools/common/container/heap"
	"sync"
	"sync/atomic"
	"time"
)

type delayQueue struct {
	h        *heap.Heap[*bucket]
	mutex    sync.Mutex
	sleeping int32         //
	wakeup   chan struct{} //唤醒通道
}

// newDelayQueue
//  @Description:  创建延迟队列
//  @return *delayQueue
//  @Author  ahKevinXy
//  @Date2023-03-29 18:21:55
func newDelayQueue() *delayQueue {
	return &delayQueue{
		h: heap.New(nil, func(e1, e2 *bucket) bool {
			return e1.getExpiration() < e2.getExpiration()
		}),

		wakeup: make(chan struct{}),
	}
}

// push
//  @Description:  添加延迟元素到队列
//  @receiver q
//  @param b
//  @Author  ahKevinXy
//  @Date2023-03-29 18:23:36
func (q *delayQueue) push(b *bucket) {
	q.mutex.Lock()

	defer q.mutex.Unlock()
	q.h.Push(b)

	if q.h.Peek() == b {
		if atomic.CompareAndSwapInt32(&q.sleeping, 1, 0) {
			q.wakeup <- struct{}{}
		}
	}
}

// take
//  @Description:  等待直到有元素到期
//  @receiver q
//  @param ctx
//  @param nowF
//  @return *bucket
//  @Author  ahKevinXy
//  @Date2023-03-29 18:24:19
func (q *delayQueue) take(ctx context.Context, nowF func() int64) *bucket {
	for {
		var t *time.Timer
		q.mutex.Lock()
		// 有元素
		if !q.h.Empty() {
			// 获取元素
			entry := q.h.Peek()
			expiration := entry.getExpiration()
			now := nowF()
			if now > expiration {
				q.h.Pop()
				q.mutex.Unlock()
				return entry
			}
			// 到期时间，使用time.NewTimer()才能够调用Stop()，从而释放定时器
			t = time.NewTimer(time.Duration(now-expiration) * time.Millisecond)
		}
		// 走到这里表示需要等待了，则需要告诉Push()在有新元素时要通知
		atomic.StoreInt32(&q.sleeping, 1)
		q.mutex.Unlock()

		// 不为空，需要同时等待元素到期，并且除非t到期，否则都需要关闭t避免泄露
		if t != nil {
			select {
			case <-q.wakeup: // 新的更快到期元素
				t.Stop()
			case <-t.C: // 首元素到期
				if atomic.SwapInt32(&q.sleeping, 0) == 0 {
					// 避免Push()的协程被阻塞
					<-q.wakeup
				}
			case <-ctx.Done(): // 被关闭
				t.Stop()
				return nil
			}
		} else {
			select {
			case <-q.wakeup: // 新的更快到期元素
			case <-ctx.Done(): // 被关闭
				return nil
			}
		}
	}
}

// channel
//  @Description:  返回一个通道，输出到期元素
//  @receiver q
//  @param ctx
//  @param size
//  @param nowF
//  @return <-chan
//  @Author  ahKevinXy
//  @Date2023-03-29 18:24:27
func (q *delayQueue) channel(ctx context.Context, size int, nowF func() int64) <-chan *bucket {
	out := make(chan *bucket, size)
	go func() {
		for {
			entry := q.take(ctx, nowF)
			if entry == nil {
				close(out)
				return
			}
			out <- entry
		}
	}()
	return out
}
