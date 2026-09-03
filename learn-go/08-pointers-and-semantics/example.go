package main

import "fmt"

type Counter struct{ Count int }

func (c Counter) IncValue()   { c.Count++ }
func (c *Counter) IncPointer() { c.Count++ }

func main() {
	c := Counter{Count: 0}
	c.IncValue()
	fmt.Println("Value recv:", c.Count)    // 0
	c.IncPointer()
	fmt.Println("Pointer recv:", c.Count)  // 1
}
