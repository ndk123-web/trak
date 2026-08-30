package ui

import "fmt"

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func Success(msg string) {
	fmt.Printf("%s✔ %s%s\n", Green, msg, Reset)
}

func Error(msg string) {
	fmt.Printf("%s✘ %s%s\n", Red, msg, Reset)
}

func Info(msg string) {
	fmt.Printf("%s➜ %s%s\n", Blue, msg, Reset)
}
