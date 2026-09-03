package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

// TODO 1: Spawn N goroutines that each increment a shared counter.
// Use sync/atomic to avoid data races. Return final counter value.
func AtomicCounter(n int) int64 {
	// FILL HERE
	return 0
}

// TODO 2: Use sync.WaitGroup to run 3 goroutines and wait for all.
// Each goroutine appends its ID to a shared slice.
// Protect the slice with sync.Mutex.
func ConcurrentAppend() []int {
	// FILL HERE
	return nil
}

// TODO 3: Return the number of logical CPUs using runtime.
func NumCPUs() int {
	// FILL HERE
	return 0
}

func main() {}
