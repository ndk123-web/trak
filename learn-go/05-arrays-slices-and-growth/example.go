package main

import "fmt"
func main() {
	s := make([]int, 2, 4)
	s = append(s, 1, 2, 3)
	fmt.Println(len(s), cap(s))
	big := make([]int, 1000)
	window := big[100:104]
	fmt.Println(len(window), cap(window))
}