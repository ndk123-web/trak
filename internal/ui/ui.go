package ui

import "fmt"

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[90m"
	White   = "\033[97m"
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
