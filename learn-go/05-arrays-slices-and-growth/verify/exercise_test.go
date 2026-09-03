package main

import "testing"

func TestGrow(t *testing.T) {
	if GrowToFive() < 5 { t.Error() }
}

func TestIndependent(t *testing.T) {
	src := []int{0,1,2,3,4,5,6,7,8,9}
	ind := IndependentSlice(src, 3, 6)
	if len(ind) != 3 || ind[0] != 3 { t.Error() }
	if SameBackingArray(src, ind) { t.Error("shares backing") }
}

func TestSameBacking(t *testing.T) {
	src := []int{1,2,3,4,5}
	sub := src[1:3]
	if !SameBackingArray(src, sub) { t.Error() }
}