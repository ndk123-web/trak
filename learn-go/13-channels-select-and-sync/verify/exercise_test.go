package main

import (
	"sort"
	"testing"
	"time"
)

func TestBufferedSum(t *testing.T) {
	if BufferedSum() != 30 { t.Error() }
}

func TestSelectWithTimeout(t *testing.T) {
	ch := make(chan string)
	go func() { time.Sleep(10 * time.Millisecond); ch <- "ok" }()
	if SelectWithTimeout(ch) != "ok" { t.Error() }

	empty := make(chan string)
	if SelectWithTimeout(empty) != "timeout" { t.Error() }
}

func TestWorkerPool(t *testing.T) {
	got := WorkerPool([]int{1, 2, 3, 4})
	if len(got) != 4 { t.Fatalf("len=%d", len(got)) }
	sort.Ints(got)
	want := []int{2, 4, 6, 8}
	for i := range want {
		if got[i] != want[i] { t.Errorf("got %v, want %v", got, want) }
	}
}
