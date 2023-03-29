package pqueue

import "github.com/ahKevinXy/go-web-tools/common/container/heap"

// PriorityQueue 优先队列
type PriorityQueue[T any] struct {
	h *heap.Heap[T]
}

func New[T any](h []T, less func(e1 T, e2 T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		h: heap.New(h, less),
	}
}

// Push 入队
func (p *PriorityQueue[T]) Push(elem T) {
	p.h.Push(elem)
}

// Pop 出队
func (p *PriorityQueue[T]) Pop() T {
	return p.h.Pop()
}

// Peek 队头元素
func (p *PriorityQueue[T]) Peek() T {
	return p.h.Peek()
}

// Len 队列元素个数
func (p *PriorityQueue[T]) Len() int {
	return p.h.Len()
}

// Empty 队列是否为空
func (p *PriorityQueue[T]) Empty() bool {
	return p.Len() == 0
}
