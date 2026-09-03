package main

import "testing"

func TestBuildInfo(t *testing.T) {
	info := BuildInfo()
	if info == "" { t.Error("empty") }
}

func TestIsLinux(t *testing.T) {
	_ = IsLinux()
}