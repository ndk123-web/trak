package main

import "testing"

func TestRectangle(t *testing.T) {
	r := Rectangle{Width: 10, Height: 5}
	if r.Area() != 50 { t.Error() }
	r.Scale(2)
	if r.Width != 20 || r.Height != 10 { t.Error() }
}

func TestSwap(t *testing.T) {
	a, b := 5, 10
	Swap(&a, &b)
	if a != 10 || b != 5 { t.Error() }
}

func TestToPointers(t *testing.T) {
	items := []Rectangle{{Width: 1, Height: 1}, {Width: 2, Height: 2}}
	ptrs := ToPointers(items)
	if len(ptrs) != 2 { t.Fatal() }
	if ptrs[0].Width != 1 || ptrs[1].Width != 2 { t.Error() }
}
