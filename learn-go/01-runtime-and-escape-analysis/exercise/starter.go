package main

// TODO 1: Return string by value (not pointer) to avoid escape
func GetName() *string {
	n := "Gopher"
	return &n
}

// TODO 2: Use concrete type to avoid interface escape
func PrintNum(n interface{}) {
	// FILL
}

// TODO 3: Pre-allocate slice to avoid backing array escape
func Grow() []int {
	data := make([]int, 2, 2)
	data = append(data, 1, 2, 3)
	return data
}

func main() {}