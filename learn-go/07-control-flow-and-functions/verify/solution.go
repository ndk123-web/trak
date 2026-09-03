package main

import "fmt"

func PrintStartEnd() {
	fmt.Println("start")
	defer fmt.Println("end")
	fmt.Println("middle")
}

func SafeIndex(slice []int, idx int) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			// We need named returns to set err from defer
		}
	}()
	if idx < 0 || idx >= len(slice) {
		return 0, fmt.Errorf("index %d out of range [0,%d)", idx, len(slice))
	}
	return slice[idx], nil
}

func DeferredSum() int {
	sum := 0
	for i := 1; i <= 3; i++ {
		defer func(n int) { sum += n }(i)
	}
	return sum
}

func main() {}
