package main

import "testing"

func TestGetName(t *testing.T) {
	name := GetName()
	if name != "Gopher" { t.Errorf("got %q", name) }
}

func TestGrow(t *testing.T) {
	s := Grow()
	if len(s) != 3 { t.Errorf("len=%d", len(s)) }
}