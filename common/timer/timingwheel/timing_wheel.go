package timingwheel

import (
	"context"
	"sync/atomic"
	"time"
	"unsafe"
)

const delayQueueBufferSize = 10 // 延迟队列缓冲区大小

// TimingWheel
//  @Description:   基于https://github.com/RussellLuo/timingwheel的实现
//  @Author  ahKevinXy
//  @Date2023-03-29 18:29:08
type TimingWheel struct {
	tick          int64          // 每一跳时间
	wheelSize     int64          //时间轮
	interval      int64          // 一圈的时间
	currentTime   int64          // 当前时间
	buckets       []*bucket      // 时间轮的每个桶
	queue         *delayQueue    //桶延迟队列
	overflowWheel unsafe.Pointer // 下一个时间轮
}

func newTimingWheel(tick, wheelSize, currentTime int64, queue *delayQueue) *TimingWheel {
	tw := &TimingWheel{
		tick:        tick,
		wheelSize:   wheelSize,
		interval:    tick * wheelSize,
		currentTime: truncate(currentTime, tick),
		buckets:     make([]*bucket, wheelSize),
		queue:       queue,
	}
	for i := 0; i < int(wheelSize); i++ {
		tw.buckets[i] = newBucket()
	}
	return tw
}
func New(tick, wheelSize int64) *TimingWheel {

	return newTimingWheel(tick, wheelSize, time.Now().UnixMilli(), newDelayQueue())
}

// Run
//  @Description:  运行时间轮
//  @receiver tw
//  @param ctx
//  @Author  ahKevinXy
//  @Date2023-03-29 18:27:09
func (tw *TimingWheel) Run(ctx context.Context) {
	bucketChan := tw.queue.channel(ctx, delayQueueBufferSize, func() int64 {
		return time.Now().UnixMilli()
	})
	for {
		select {
		case b := <-bucketChan: // 桶到期
			// 前进当前时间
			tw.advance(b.expiration)
			// 处理桶
			b.flush(tw.addOrRun)
		case <-ctx.Done(): // 被关闭
			return
		}
	}
}

// AfterFunc
//  @Description:  添加定时器
//  @receiver tw
//  @param delay
//  @param f
//  @return *Timer
//  @Author  ahKevinXy
//  @Date2023-03-29 18:27:22
func (tw *TimingWheel) AfterFunc(delay time.Duration, f func()) *Timer {
	t := &Timer{
		expiration: time.Now().Add(delay).UnixMilli(),
		task:       f,
	}
	tw.add(t)
	return t
}

type Scheduler interface {
	// Next 表示下一个执行任务的时间
	// 如果time.IsZero()==true则不再进行
	Next(time.Time) time.Time
}

func (tw *TimingWheel) ScheduleFunc(s Scheduler, f func()) (t *Timer) {
	expiration := s.Next(time.Now())
	if expiration.IsZero() {
		return
	}

	t = &Timer{
		expiration: expiration.UnixMilli(),
		task: func() {
			// 添加下一次执行任务
			expiration := s.Next(time.UnixMilli(t.expiration))
			if !expiration.IsZero() {
				t.expiration = expiration.UnixMilli()
				tw.addOrRun(t)
			}
			// 执行任务
			f()
		},
	}
	tw.addOrRun(t)
	return
}

// add
//  @Description:  添加定时器
//  @receiver tw
//  @param t
//  @return bool
//  @Author  ahKevinXy
//  @Date2023-03-29 18:28:05
func (tw *TimingWheel) add(t *Timer) bool {
	currentTime := atomic.LoadInt64(&tw.currentTime)
	if t.expiration < currentTime+tw.tick { // 已经过期了
		return false
	} else if t.expiration < currentTime+tw.interval { // 在当前时间轮里
		// 多少跳了
		ticks := t.expiration / tw.tick
		// 应该在时间轮的哪个桶里
		b := tw.buckets[ticks%tw.wheelSize]
		b.add(t)

		// 如果设置桶过期时间成功
		// 表示这个桶第一次加入定时器，因此应该把它放到延迟队列里面去等待到期
		if b.setExpiration(ticks * tw.tick) {
			tw.queue.push(b)
		}
		return true
	} else { // 在其他时间轮里
		overflowWheel := atomic.LoadPointer(&tw.overflowWheel)
		if overflowWheel == nil {
			tw.setOverflowWheel(currentTime)
			overflowWheel = atomic.LoadPointer(&tw.overflowWheel)
		}
		return (*TimingWheel)(overflowWheel).add(t)
	}
}

// addOrRun
//  @Description:  添加任务或运行
//  @receiver tw
//  @param t
//  @Author  ahKevinXy
//  @Date2023-03-29 18:28:25
func (tw *TimingWheel) addOrRun(t *Timer) {
	if !tw.add(t) {
		go t.task()
	}
}

// advance
//  @Description:  前进时间
//  @receiver tw
//  @param expiration
//  @Author  ahKevinXy
//  @Date2023-03-29 18:28:40
func (tw *TimingWheel) advance(expiration int64) {
	currentTime := atomic.LoadInt64(&tw.currentTime)
	if expiration >= currentTime+tw.tick {
		currentTime := truncate(expiration, tw.tick)
		atomic.StoreInt64(&tw.currentTime, currentTime)

		overflowWheel := atomic.LoadPointer(&tw.overflowWheel)
		if overflowWheel != nil {
			(*TimingWheel)(overflowWheel).advance(currentTime)
		}
	}
}

func (tw *TimingWheel) setOverflowWheel(currentTime int64) {
	overflowWheel := newTimingWheel(tw.interval, tw.wheelSize, currentTime, tw.queue)
	atomic.CompareAndSwapPointer(&tw.overflowWheel, nil, unsafe.Pointer(overflowWheel))
}

// truncate
//  @Description:  去除不满一整跳的时间
//  @param time
//  @param tick
//  @return int64
//  @Author  ahKevinXy
//  @Date2023-03-29 18:28:53
func truncate(time, tick int64) int64 {
	return time - time%tick
}
