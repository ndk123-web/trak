package main

import (
	"fmt"
	"os"
	"github.com/ndk123-web/trak/internal/app"
)

func main() {
	arguments := os.Args

	_, err := app.Run(arguments)
	if err != nil {
		fmt.Printf("%v", err.Error())
	}
}
