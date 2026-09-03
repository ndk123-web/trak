package main

import "encoding/json"

type Name struct {
	First string `json:"first_name"`
	Last  string `json:"last_name"`
}

type Person struct {
	Name
	Age int `json:"age"`
}

func SafeMapRead(m map[string]int, key string) (int, bool) {
	val, ok := m[key]
	return val, ok
}

func PersonToJSON(p Person) ([]byte, error) {
	return json.Marshal(p)
}

func main() {}
