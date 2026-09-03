package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCopyAll(t *testing.T) {
	r := strings.NewReader("hello")
	var w bytes.Buffer
	n, err := CopyAll(r, &w)
	if err != nil { t.Fatal(err) }
	if n != 5 || w.String() != "hello" { t.Error() }
}

func TestReadExactly(t *testing.T) {
	r := strings.NewReader("hello world")
	s, err := ReadExactly(r, 5)
	if err != nil { t.Fatal(err) }
	if s != "hello" { t.Errorf("got %q", s) }
}

func TestChainReaders(t *testing.T) {
	r1 := strings.NewReader("abc")
	r2 := strings.NewReader("def")
	ch := ChainReaders(r1, r2)
	var w bytes.Buffer
	io.Copy(&w, ch)
	if w.String() != "abcdef" { t.Errorf("got %q", w.String()) }
}
