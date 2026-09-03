package main

import (
	"strings"
	"testing"
)

func TestSafeIndex(t *testing.T) {
	s := []int{10, 20, 30}
	v, err := SafeIndex(s, 1)
	if err != nil || v != 20 { t.Error() }
	v, err = SafeIndex(s, 10)
	if err == nil { t.Error("expected error") }
	if !strings.Contains(err.Error(), "index") && !strings.Contains(err.Error(), "range") {
		t.Errorf("bad error msg: %v", err)
	}
}

func TestDeferredSum(t *testing.T) {
	if DeferredSum() != 6 { t.Error() }
}
