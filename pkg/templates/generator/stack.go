package generator

type Stack[T any] struct {
	stack []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		stack: make([]T, 0),
	}
}

func (s *Stack[T]) Push(value T) {
	s.stack = append(s.stack, value)
}

func (s *Stack[T]) Pop() (result T) {
	if s.IsEmpty() {
		return
	}
	index := len(s.stack) - 1
	element := s.stack[index]
	s.stack = s.stack[:index]
	return element
}

func (s *Stack[T]) Peek() (result T) {
	if s.IsEmpty() {
		return
	}
	return s.stack[len(s.stack)-1]
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.stack) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.stack)
}
