package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

func QueryDB() error {
	return fmt.Errorf("db query: %w", ErrNotFound)
}

func main() {
	err := QueryDB()
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Detected ErrNotFound in chain!")
	}
}
