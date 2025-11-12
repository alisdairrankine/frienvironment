package frienvironment

//--stack

type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Pop() T {
	n := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return n
}

func (s *Stack[T]) Push(n T) {
	s.data = append(s.data, n)
}

func (s *Stack[T]) PushMany(n ...T) {
	s.data = append(s.data, n...)
}
