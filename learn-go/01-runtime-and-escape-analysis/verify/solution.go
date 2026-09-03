package main

import "fmt"

func GetName() string {
	n := "Gopher"
	return n
}

func PrintNum(n int) {
	fmt.Println(n)
}

func Grow() []int {
	data := make([]int, 0, 5)
	data = append(data, 1, 2, 3)
	return data
}

func main() {}