package app

import "fmt"

func Run(arguments []string) (int, error) {

	fmt.Printf("Arguments Are:\n")
	for idx , val := range arguments {
		fmt.Print("idx: ", idx, " ")
		fmt.Print("Val: ", val, "\n")
	}

	return 1, nil
}