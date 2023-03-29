package deque

import "github.com/ahKevinXy/go-web-tools/common/container/list"

type Stack[T any] struct {
	l *list.List[T]
}

func New[T any]() *Stack[T] {
	return &Stack[T]{l: list.New[T]()}
}

// Push 入栈
func (s *Stack[T]) Push(elem T) {
	s.l.PushBack(elem)
}

// Pop 出栈
func (s *Stack[T]) Pop() T {
	return s.l.RemoveBack()
}

// Peek 栈顶元素
func (s *Stack[T]) Peek() T {
	return s.l.Back().Value
}

// Len 栈元素个数
func (s *Stack[T]) Len() int {
	return s.l.Len()
}

// Empty 栈是否为空
func (s *Stack[T]) Empty() bool {
	return s.l.Empty()
}
