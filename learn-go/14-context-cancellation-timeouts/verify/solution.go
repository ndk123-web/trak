package main

import (
	"context"
	"time"
)

func CheckTimeout() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	select {
	case <-time.After(50 * time.Millisecond):
		return false
	case <-ctx.Done():
		return true
	}
}

func ParentCancelsChild() bool {
	parent, cancel := context.WithCancel(context.Background())
	child, childCancel := context.WithTimeout(parent, time.Minute)
	defer childCancel()
	cancel()
	select {
	case <-child.Done():
		return true
	case <-time.After(100 * time.Millisecond):
		return false
	}
}

type requestIDKey struct{}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}

func main() {}
