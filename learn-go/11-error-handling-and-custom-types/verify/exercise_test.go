package main

import (
	"errors"
	"testing"
)

func TestValidate(t *testing.T) {
	err := Validate(-5)
	if err == nil { t.Fatal("expected error") }
	if !errors.Is(err, ErrInvalidInput) { t.Error("should wrap ErrInvalidInput") }
}

func TestCheckValidation(t *testing.T) {
	ve := &ValidationError{Field: "age", Message: "too young"}
	wrapped := fmt.Errorf("wrapped: %w", ve)
	found, ok := CheckValidation(wrapped)
	if !ok { t.Fatal("expected to find ValidationError") }
	if found.Field != "age" { t.Error() }
}
