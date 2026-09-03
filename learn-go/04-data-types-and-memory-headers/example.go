package main

import ("fmt"; "unsafe"; "unicode/utf8")
func main() {
	fmt.Println(unsafe.Sizeof(""), unsafe.Sizeof([]int{}))
	s := "Hello, 世界"
	fmt.Println(len(s), utf8.RuneCountInString(s))
}