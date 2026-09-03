package main

import (
	"context"
	"time"
)

// TODO 1: Implement a worker pool that squares integers.
// Spawn 'workers' goroutines. Each reads from jobs chan, sends square to results chan.
// Return all results as a slice. Use context with 100ms timeout.
func SquareWorkerPool(inputs []int, workers int) []int {
	// FILL HERE
	return nil
}

// TODO 2: Fix the nil interface bug. When ok=false, return nil error.
// Currently returns typed nil pointer inside interface.
type AppError struct{ Msg string }
func (e *AppError) Error() string { return e.Msg }

func MaybeError(ok bool) error {
	var err *AppError = nil
	if !ok {
		return err // BUG
	}
	return nil
}

// TODO 3: Write a function that demonstrates defer LIFO order.
// It should return the string "321" by appending digits in defer blocks.
func DeferLIFO() string {
	// FILL HERE
	return ""
}

func main() {}
