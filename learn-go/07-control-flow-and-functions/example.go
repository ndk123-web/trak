package main

import "fmt"

func SafeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()
	result = a / b
	return
}

func main() {
	for i := 1; i <= 3; i++ {
		defer fmt.Printf("defer %d\n", i)
	}
	res, err := SafeDivide(10, 0)
	fmt.Printf("Result: %d, Error: %v\n", res, err)
}
