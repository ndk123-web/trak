package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	select {
	case <-time.After(100 * time.Millisecond):
		fmt.Println("Finished")
	case <-ctx.Done():
		fmt.Println("Timeout:", ctx.Err())
	}
}
