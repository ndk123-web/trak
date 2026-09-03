package main

import (
	"context"
	"fmt"
	"time"
)

func SquareWorkerPool(inputs []int, workers int) []int {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	jobs := make(chan int, len(inputs))
	results := make(chan int, len(inputs))
	for w := 0; w < workers; w++ {
		go func() {
			for {
				select {
				case j, ok := <-jobs:
					if !ok {
						return
					}
					select {
					case results <- j * j:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	for _, v := range inputs {
		jobs <- v
	}
	close(jobs)
	out := make([]int, 0)
	for i := 0; i < len(inputs); i++ {
		select {
		case r := <-results:
			out = append(out, r)
		case <-ctx.Done():
			return out
		}
	}
	return out
}

func MaybeError(ok bool) error {
	if !ok {
		return nil
	}
	return &AppError{Msg: "error"}
}

func DeferLIFO() string {
	var s string
	for i := 1; i <= 3; i++ {
		defer func(n int) { s += fmt.Sprintf("%d", n) }(i)
	}
	return s
}

type AppError struct{ Msg string }
func (e *AppError) Error() string { return e.Msg }

func main() {}
