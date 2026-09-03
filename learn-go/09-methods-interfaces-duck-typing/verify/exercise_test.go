package main

import "testing"

func TestGetError(t *testing.T) {
	if GetError(false) != nil { t.Error("should be nil") }
}

func TestTypeName(t *testing.T) {
	if TypeName(42) != "int" { t.Error() }
	if TypeName("hello") != "string" { t.Error() }
	if TypeName(true) != "bool" { t.Error() }
}

func TestIsTrulyNil(t *testing.T) {
	if !IsTrulyNil(nil) { t.Error() }
	var e *CustomError = nil
	if IsTrulyNil(e) { t.Error("typed nil should not be truly nil") }
}
