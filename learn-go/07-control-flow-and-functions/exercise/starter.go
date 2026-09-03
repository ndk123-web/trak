package main

import "fmt"

// TODO 1: Use defer to print "start" then "end" regardless of what happens.
func PrintStartEnd() {
	// FILL HERE
	fmt.Println("middle")
}

// TODO 2: Recover from slice index out of bounds panic. Return error instead.
func SafeIndex(slice []int, idx int) (int, error) {
	// FILL HERE
	return 0, nil
}

// TODO 3: Return sum of all deferred increments (1+2+3=6).
// Use closure and defer. LIFO means 3+2+1.
func DeferredSum() int {
	// FILL HERE
	return 0
}

func main() {}
