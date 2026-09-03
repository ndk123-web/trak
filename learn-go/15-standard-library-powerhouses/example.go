package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func main() {
	r := strings.NewReader("pipeline data")
	var w bytes.Buffer
	io.Copy(&w, r)
	fmt.Println(w.String())
}
