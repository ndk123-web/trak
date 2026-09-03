package main

import (
	"io"
)

func CopyAll(r io.Reader, w io.Writer) (int64, error) {
	return io.Copy(w, r)
}

func ReadExactly(r io.Reader, n int) (string, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	return string(buf), err
}

func ChainReaders(r1, r2 io.Reader) io.Reader {
	return io.MultiReader(r1, r2)
}

func main() {}
