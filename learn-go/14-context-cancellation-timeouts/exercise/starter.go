package main

import (
	"context"
	"time"
)

// TODO 1: Create a context with 20ms timeout. Simulate work that takes 50ms.
// Return true if timeout occurred, false if work finished.
func CheckTimeout() bool {
	// FILL HERE
	return false
}

// TODO 2: Create a parent context with cancel. Create a child with timeout.
// Cancel the PARENT and return true if the child context is also cancelled.
func ParentCancelsChild() bool {
	// FILL HERE
	return false
}

// TODO 3: Use context.WithValue to carry a request ID.
// Define a private key type. Return the request ID from context.
type requestIDKey struct{}

func GetRequestID(ctx context.Context) string {
	// FILL HERE
	return ""
}

func main() {}
