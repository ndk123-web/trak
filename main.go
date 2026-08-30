package main

import (
	"fmt"

	"github.com/ndk123-web/trak/cmd"
	"github.com/ndk123-web/trak/internal/config"
)

func main() {
	config.UpdateBaseUrl()

	if err := cmd.Execute(); err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
