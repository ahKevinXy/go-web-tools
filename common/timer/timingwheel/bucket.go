package timingwheel

import (
	"github.com/ahKevinXy/go-web-tools/common/container/list"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Timer struct {
	expiration int64                 // 到期时间
	task       func()                // 任务
	b          unsafe.Pointer        //所属时间轮的桶
	elem       *list.Element[*Timer] //为了能从链表中删除
}

func (t *Timer) Stop() bool {
	stop := false

	for b := t.getBucket(); b != nil; b = t.getBucket() {
		stop = b.remove(t)
	}

	return stop
}

func (t *Timer) getBucket() *bucket {
	return (*bucket)(atomic.LoadPointer(&t.b))
}

func (t *Timer) setBucket(b *bucket) {
	atomic.StorePointer(&t.b, unsafe.Pointer(b))
}

type bucket struct {
	expiration int64              // 到期时间
	timers     *list.List[*Timer] // 定时器;列表
	mutex      sync.Mutex         // 并发锁
}

func (b *bucket) add(t *Timer) {
	b.mutex.Lock()

	defer b.mutex.Unlock()
	elem := b.timers.PushBack(t)
	t.elem = elem
	t.setBucket(b)
}

func (b *bucket) remove(t *Timer) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if t.getBucket() != b {
		return false
	}

	b.timers.Remove(t.elem)
	t.setBucket(nil)
	t.elem = nil
	return true
}

func newBucket() *bucket {
	return &bucket{
		expiration: -1,
		timers:     list.New[*Timer](),
	}
}

// 添加到上一级定时器或执行任务
func (b *bucket) flush(addOrRun func(t *Timer)) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for elem := b.timers.Front(); elem != nil; {
		next := elem.Next()
		t := elem.Value
		if t.getBucket() == b {
			t.setBucket(nil)
			t.elem = nil
		}
		addOrRun(t)
		elem = next
	}

	// 设置过期时间表示没有加入到延迟队列
	b.setExpiration(-1)
	b.timers.Clear()
}

func (b *bucket) getExpiration() int64 {
	return atomic.LoadInt64(&b.expiration)
}

// 返回true表示设置成功
// 否则表示没变化
func (b *bucket) setExpiration(expiration int64) bool {
	return atomic.SwapInt64(&b.expiration, expiration) != expiration
}
