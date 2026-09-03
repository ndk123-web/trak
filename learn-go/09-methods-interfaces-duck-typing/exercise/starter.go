package main

import "fmt"

// TODO 1: Fix GetError. When ok=false, returned error MUST be nil.
type CustomError struct{ Msg string }
func (e *CustomError) Error() string { return e.Msg }

func GetError(ok bool) error {
	var err *CustomError = nil
	if !ok {
		return err // BUG!
	}
	return nil
}

// TODO 2: Return underlying type name using type switch.
func TypeName(v any) string {
	// FILL HERE
	return ""
}

// TODO 3: Return true ONLY if error is truly nil (both type and value nil).
func IsTrulyNil(err error) bool {
	// FILL HERE
	return false
}

func main() {}
