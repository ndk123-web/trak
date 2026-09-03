package main

import "fmt"

type MyError struct{ Code int }
func (e *MyError) Error() string { return fmt.Sprintf("code: %d", e.Code) }

func Buggy(ok bool) error {
	var err *MyError = nil
	if !ok { return err }
	return nil
}

func Fixed(ok bool) error {
	if !ok { return nil }
	return nil
}

func main() {
	fmt.Println("Buggy != nil?", Buggy(false) != nil)
	fmt.Println("Fixed != nil?", Fixed(false) != nil)
}
