package main

import "fmt"

type Number interface {
	~int | ~float64
}

func Add[T Number](a, b T) T { return a + b }

func Contains[T comparable](s []T, t T) bool {
	for _, v := range s { if v == t { return true } }
	return false
}

func main() {
	fmt.Println(Add(10, 20))
	fmt.Println(Contains([]string{"a", "b"}, "b"))
}
