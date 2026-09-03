package main

import (
	"errors"
	"fmt"
)

var ErrInvalidInput = errors.New("invalid input")

// TODO 1: Wrap ErrInvalidInput with context "validation failed".
// Use %w so errors.Is works.
func Validate(age int) error {
	if age < 0 {
		// FILL HERE
		return nil
	}
	return nil
}

// TODO 2: Define a custom error type ValidationError with Field and Message.
// Implement the Error() string method.
type ValidationError struct {
	// FILL HERE
}

// TODO 3: Write a function CheckValidation that uses errors.As to extract
// a *ValidationError from an error chain. Return it and true if found.
func CheckValidation(err error) (*ValidationError, bool) {
	// FILL HERE
	return nil, false
}

func main() {}
