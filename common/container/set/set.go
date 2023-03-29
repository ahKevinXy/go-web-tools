package set

type Set[T comparable] struct {
	m map[T]struct{}
}

func New[T comparable]() *Set[T] {
	return &Set[T]{m: map[T]struct{}{}}
}

// Add 加入集合
func (s *Set[T]) Add(elem T) {
	s.m[elem] = struct{}{}
}

// Remove 移出集合
func (s *Set[T]) Remove(elem T) {
	delete(s.m, elem)
}

// Contains 是否包含元素
func (s *Set[T]) Contains(elem T) bool {
	_, contains := s.m[elem]
	return contains
}

// Len 集合长度
func (s *Set[T]) Len() int {
	return len(s.m)
}

// Empty 集合是否为空
func (s *Set[T]) Empty() bool {
	return s.Len() == 0
}
