package main

import "fmt"

// Stack is a generic LIFO structure.
// T can be any type.
type Stack[T any] struct {
	items []T
}

// Push adds an element to the stack
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop removes and returns the last element
func (s *Stack[T]) Pop() T {
	last := len(s.items) - 1
	val := s.items[last]
	s.items = s.items[:last]
	return val
}

func main() {
	var s Stack[int]

	s.Push(10)
	s.Push(20)
	s.Push(30)

	fmt.Println(s.Pop()) // 30
	fmt.Println(s.Pop()) // 20
}
