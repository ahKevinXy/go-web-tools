package delayqueue

import (
	"context"
	"github.com/ahKevinXy/go-web-tools/common/container/heap"
	"sync"
	"sync/atomic"
	"time"
)

// 延时队列

type entry[T any] struct {
	value      T
	expiration time.Time //到期时间
}

type DelayQueue[T any] struct {
	h     *heap.Heap[*entry[T]]
	mutex sync.Mutex //保证并发安全
	// 0 表示take 没在等待 1 表示在等待
	sleeping int32
	// 唤醒通道
	wakeup chan struct{}
}

// New [T any]
//  @Description: 创建延时队列
//  @return *DelayQueue[T]
func New[T any]() *DelayQueue[T] {
	return &DelayQueue[T]{
		h: heap.New(nil, func(e1, e2 *entry[T]) bool {
			return e1.expiration.Before(e2.expiration)
		}),
		wakeup: make(chan struct{}),
	}
}

// Push
//  @Description: 添加延迟元素到队列
//  @receiver q
//  @param value
//  @param delay
//  @Author  ahKevinXy
//  @Date2023-03-29 17:46:03
func (q *DelayQueue[T]) Push(value T, delay time.Duration) {
	// 并发锁
	q.mutex.Lock()

	defer q.mutex.Unlock()
	entry := &entry[T]{
		value:      value,
		expiration: time.Now().Add(delay),
	}

	q.h.Push(entry)

	// 唤醒等待的take

	if q.h.Peek() == entry {
		// 把sleep 从1修改成0
		if atomic.CompareAndSwapInt32(&q.sleeping, 1, 0) {
			q.wakeup <- struct{}{}
		}
	}

}

// Take
//  @Description: 等待直到有元素到期
//  @receiver q
//  @param ctx
//  @return T
//  @return bool
//  @Author  ahKevinXy
//  @Date2023-03-29 17:53:38
func (q *DelayQueue[T]) Take(ctx context.Context) (T, bool) {
	for {
		var timer *time.Timer
		q.mutex.Lock()
		// 判断是否为空
		if !q.h.Empty() {
			entry := q.h.Peek()
			now := time.Now()
			if now.After(entry.expiration) {
				q.h.Pop()
				q.mutex.Unlock()
				return entry.value, true
			}
			// 到期时间,使用 time.NewAfter 才能调用stop

			timer = time.NewTimer(entry.expiration.Sub(now))
		}

		atomic.StoreInt32(&q.sleeping, 1)
		q.mutex.Unlock()

		if timer != nil {
			select {
			case <-q.wakeup:
				timer.Stop()
			case <-timer.C:
				if atomic.SwapInt32(&q.sleeping, 0) == 0 {
					<-q.wakeup
				}
			case <-ctx.Done():
				timer.Stop()
				var t T
				return t, false

			}
		} else {
			select {
			case <-q.wakeup: //
			case <-ctx.Done(): // 被关闭
				var t T
				return t, false
			}
		}
	}
}

// Channel
//  @Description: 返回一个通道，输出到期元素
//  @receiver q
//  @param ctx
//  @param size
//  @return <-chan
//  @Author  ahKevinXy
//  @Date2023-03-29 17:55:47
func (q *DelayQueue[T]) Channel(ctx context.Context, size int) <-chan T {
	out := make(chan T, size)
	go func() {
		for {
			entry, ok := q.Take(ctx)
			if !ok {
				close(out)
				return
			}
			out <- entry
		}
	}()

	return out
}

// Peek
//  @Description: 获取队头元素
//  @receiver q
//  @return T
//  @return bool
//  @Author  ahKevinXy
//  @Date2023-03-29 17:57:34
func (q *DelayQueue[T]) Peek() (T, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.h.Empty() {
		var t T
		return t, false
	}
	return q.h.Peek().value, true
}

// Pop
//  @Description: 获取到期元素
//  @receiver q
//  @return T
//  @return bool
//  @Author  ahKevinXy
//  @Date2023-03-29 17:58:34
func (q *DelayQueue[T]) Pop() (T, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	// 没元素
	if q.h.Empty() {
		var t T
		return t, false
	}
	entry := q.h.Peek()
	// 还没元素到期
	if time.Now().Before(entry.expiration) {
		var t T
		return t, false
	}
	// 移除元素
	q.h.Pop()
	return entry.value, true
}

// Empty
//  @Description: 是否队列为空
//  @receiver q
//  @return bool
//  @Author  ahKevinXy
//  @Date2023-03-29 17:58:56
func (q *DelayQueue[T]) Empty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.h.Empty()
}
