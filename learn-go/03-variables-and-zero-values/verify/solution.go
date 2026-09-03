package main

type LogLevel int
const (Debug LogLevel = iota; Info; Warn; Error)
const (PermCreate = 1 << iota; PermRead; PermUpdate; PermDelete)
func ZeroValues() map[string]interface{} {
	return map[string]interface{}{"int":0, "string":"", "bool":false, "slice":[]int(nil), "map":map[string]int(nil), "ptr":(*int)(nil)}
}
func main() {}