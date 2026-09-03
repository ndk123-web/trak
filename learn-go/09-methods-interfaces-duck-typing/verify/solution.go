package main

import "fmt"

type CustomError struct{ Msg string }
func (e *CustomError) Error() string { return e.Msg }

func GetError(ok bool) error {
	if !ok {
		return nil
	}
	return &CustomError{Msg: "error"}
}

func TypeName(v any) string {
	switch v.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	default:
		return fmt.Sprintf("%T", v)
	}
}

func IsTrulyNil(err error) bool {
	return err == nil
}

func main() {}
