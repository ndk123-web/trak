package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	fmt.Printf("Logical CPUs: %d\n", runtime.GOMAXPROCS(0))
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d\n", id)
		}(i)
	}
	wg.Wait()
}
