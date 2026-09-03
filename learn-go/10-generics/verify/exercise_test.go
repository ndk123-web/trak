package main

import "testing"

func TestMax(t *testing.T) {
	if Max(5, 10) != 10 { t.Error() }
	if Max(3.14, 2.71) != 3.14 { t.Error() }
}

func TestFilter(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	if len(evens) != 2 || evens[0] != 2 || evens[1] != 4 { t.Error() }
}

func TestStack(t *testing.T) {
	var s Stack[int]
	s.Push(10)
	s.Push(20)
	v, ok := s.Pop()
	if !ok || v != 20 { t.Error() }
	v, ok = s.Pop()
	if !ok || v != 10 { t.Error() }
	_, ok = s.Pop()
	if ok { t.Error() }
}
