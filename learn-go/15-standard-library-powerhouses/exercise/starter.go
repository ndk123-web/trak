package main

import (
	"bytes"
	"io"
	"strings"
)

// TODO 1: Copy all data from reader to writer using io.Copy.
// Return the number of bytes copied.
func CopyAll(r io.Reader, w io.Writer) (int64, error) {
	// FILL HERE
	return 0, nil
}

// TODO 2: Read exactly N bytes from a reader. Return them as a string.
// If reader has fewer than N bytes, return what you got and an error.
func ReadExactly(r io.Reader, n int) (string, error) {
	// FILL HERE
	return "", nil
}

// TODO 3: Chain two readers: first read from r1, then from r2.
// Return a single io.Reader that reads r1 then r2.
func ChainReaders(r1, r2 io.Reader) io.Reader {
	// FILL HERE
	return nil
}

func main() {}
