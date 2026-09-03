package main

import (
	"sort"
	"testing"
)

func TestSquareWorkerPool(t *testing.T) {
	got := SquareWorkerPool([]int{1, 2, 3, 4, 5}, 2)
	if len(got) == 0 { t.Fatal("no results") }
	sort.Ints(got)
	// At least some results should come back
	if got[0] != 1 { t.Errorf("first result = %d, want 1", got[0]) }
}

func TestMaybeError(t *testing.T) {
	if MaybeError(false) != nil { t.Error("should be nil") }
}

func TestDeferLIFO(t *testing.T) {
	if DeferLIFO() != "321" { t.Errorf("got %q", DeferLIFO()) }
}
