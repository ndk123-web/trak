package main

import (
	"fmt"
	"time"
)

// TODO 1: Create a buffered channel of ints with capacity 2.
// Send two values without blocking, then receive both and return their sum.
func BufferedSum() int {
	// FILL HERE
	return 0
}

// TODO 2: Use select with a timeout to read from a channel.
// If the channel has no value within 50ms, return "timeout".
// If it has a value, return that value as string.
func SelectWithTimeout(ch chan string) string {
	// FILL HERE
	return ""
}

// TODO 3: Implement a worker pool pattern.
// Spawn 2 goroutines that read jobs from a jobs channel and send results to a results channel.
// Each worker doubles the input. Return all results as a slice.
func WorkerPool(inputs []int) []int {
	// FILL HERE
	return nil
}

func main() {}
