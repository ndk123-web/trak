package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan string, 1)
	ch <- "msg"
	select {
	case msg := <-ch:
		fmt.Println("Received:", msg)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("Timeout")
	}

	var mu sync.RWMutex
	mu.Lock()
	fmt.Println("Write locked")
	mu.Unlock()
}
