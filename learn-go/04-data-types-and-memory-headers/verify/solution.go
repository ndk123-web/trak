package main

import "unsafe"

func StringHeaderSize() uintptr { return unsafe.Sizeof("") }
func SliceHeaderSize() uintptr { return unsafe.Sizeof([]int{}) }
func RuneCount(s string) int { count := 0; for range s { count++ }; return count }
func main() {}