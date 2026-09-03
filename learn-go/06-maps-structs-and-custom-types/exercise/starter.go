package main

import "encoding/json"

// TODO 1: Safely read from map. Return (value, true) if exists, else (0, false).
func SafeMapRead(m map[string]int, key string) (int, bool) {
	// FILL HERE
	return 0, false
}

// TODO 2: Define Person with embedded Name struct (First, Last fields).
// JSON tags: first_name, last_name, age.
type Name struct {
	// FILL HERE
}

type Person struct {
	// FILL HERE
}

// TODO 3: Marshal Person to JSON bytes.
func PersonToJSON(p Person) ([]byte, error) {
	// FILL HERE
	return nil, nil
}

func main() {}
