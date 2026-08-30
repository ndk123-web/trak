package main

import (
	"fmt"
	"github.com/ndk123-web/trak/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
