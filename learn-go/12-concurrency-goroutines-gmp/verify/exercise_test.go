package main

import (
	"sort"
	"testing"
)

func TestAtomicCounter(t *testing.T) {
	got := AtomicCounter(1000)
	if got != 1000 { t.Errorf("got %d, want 1000", got) }
}

func TestConcurrentAppend(t *testing.T) {
	got := ConcurrentAppend()
	if len(got) != 3 { t.Fatalf("len=%d", len(got)) }
	sort.Ints(got)
	if got[0] != 1 || got[1] != 2 || got[2] != 3 { t.Errorf("got %v", got) }
}

func TestNumCPUs(t *testing.T) {
	if NumCPUs() <= 0 { t.Error() }
}
