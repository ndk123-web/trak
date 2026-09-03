package main

import (
	"errors"
	"fmt"
)

var ErrInvalidInput = errors.New("invalid input")

func Validate(age int) error {
	if age < 0 {
		return fmt.Errorf("validation failed: %w", ErrInvalidInput)
	}
	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

func CheckValidation(err error) (*ValidationError, bool) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve, true
	}
	return nil, false
}

func main() {}
