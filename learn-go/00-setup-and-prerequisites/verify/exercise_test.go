package main

import "testing"

func TestGoVersion(t *testing.T) {
	v := GoVersion()
	if v == "" { t.Error("empty version") }
}

func TestGetCPUCount(t *testing.T) {
	if GetCPUCount() <= 0 { t.Error("invalid cpu count") }
}