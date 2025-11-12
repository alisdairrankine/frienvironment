package frienvironment

import "sync"

//--stack

type Stack[T any] struct {
	sync.Mutex
	data []T
}

func (s *Stack[T]) Pop() T {
	s.Lock()
	defer s.Unlock()
	n := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return n
}

func (s *Stack[T]) Push(n T) {
	s.Lock()
	defer s.Unlock()
	s.data = append(s.data, n)
}

func (s *Stack[T]) PushMany(n ...T) {
	s.Lock()
	defer s.Unlock()
	s.data = append(s.data, n...)
}
