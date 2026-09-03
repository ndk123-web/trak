package main

func GrowToFive() int {
	s := make([]int, 0, 2)
	for i := 1; i <= 5; i++ { s = append(s, i) }
	return cap(s)
}
func IndependentSlice(src []int, start, end int) []int {
	sub := src[start:end]
	out := make([]int, len(sub))
	copy(out, sub)
	return out
}
func SameBackingArray(a, b []int) bool {
	if len(a) == 0 || len(b) == 0 { return false }
	return &a[0] == &b[0]
}
func main() {}