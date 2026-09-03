package main

import "time"

func BufferedSum() int {
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	return <-ch + <-ch
}

func SelectWithTimeout(ch chan string) string {
	select {
	case v := <-ch:
		return v
	case <-time.After(50 * time.Millisecond):
		return "timeout"
	}
}

func WorkerPool(inputs []int) []int {
	jobs := make(chan int, len(inputs))
	results := make(chan int, len(inputs))
	for w := 0; w < 2; w++ {
		go func() {
			for j := range jobs {
				results <- j * 2
			}
		}()
	}
	for _, v := range inputs {
		jobs <- v
	}
	close(jobs)
	out := make([]int, 0, len(inputs))
	for i := 0; i < len(inputs); i++ {
		out = append(out, <-results)
	}
	return out
}

func main() {}
