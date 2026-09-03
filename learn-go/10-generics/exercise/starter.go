package main

// TODO 1: Generic Max with ordering constraint.
type Ordered interface {
	// FILL HERE
}

func Max[T Ordered](a, b T) T {
	// FILL HERE
	return a
}

// TODO 2: Generic Filter with predicate.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	// FILL HERE
	return nil
}

// TODO 3: Generic Stack with Push and Pop.
type Stack[T any] struct {
	// FILL HERE
}

func (s *Stack[T]) Push(v T) {
	// FILL HERE
}

func (s *Stack[T]) Pop() (T, bool) {
	// FILL HERE
	var zero T
	return zero, false
}

func main() {}
