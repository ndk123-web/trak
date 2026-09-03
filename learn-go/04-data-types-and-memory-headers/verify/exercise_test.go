package main

import "testing"

func TestStringHeader(t *testing.T) {
	if StringHeaderSize() != 16 { t.Error() }
}

func TestSliceHeader(t *testing.T) {
	if SliceHeaderSize() != 24 { t.Error() }
}

func TestRuneCount(t *testing.T) {
	if RuneCount("Hello, 世界") != 9 { t.Error() }
}