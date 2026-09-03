package main

import "fmt"

type C struct { V int }
func Heap() *C { c := C{V: 1}; return &c }
func Stack() C  { c := C{V: 2}; return c }
func main() {
	fmt.Println(Heap(), Stack())
}