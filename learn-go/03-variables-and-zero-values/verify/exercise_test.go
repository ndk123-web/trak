package main

import "testing"

func TestLogLevel(t *testing.T) {
	if Debug != 0 { t.Error() }
	if Error != 3 { t.Error() }
}

func TestPerms(t *testing.T) {
	if PermCreate != 1 || PermDelete != 8 { t.Error() }
}

func TestZero(t *testing.T) {
	z := ZeroValues()
	if z["int"] != 0 || z["bool"] != false { t.Error() }
}