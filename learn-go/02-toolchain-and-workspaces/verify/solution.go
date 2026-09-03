package main

import ("fmt"; "runtime")

func BuildInfo() string {
	return fmt.Sprintf("%s/%s %s", runtime.GOOS, runtime.GOARCH, runtime.Version())
}
func IsLinux() bool { return runtime.GOOS == "linux" }
func main() {}