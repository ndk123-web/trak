package main

import "fmt"

type Status int
const (
	Pending Status = iota
	Active
)
const (
	Read = 1 << iota
	Write
)
func main() {
	fmt.Println(Pending, Active, Read, Write)
}