package main

import "runtime"

func GoVersion() string { return runtime.Version() }
func GetCPUCount() int { return runtime.NumCPU() }
func main() {}