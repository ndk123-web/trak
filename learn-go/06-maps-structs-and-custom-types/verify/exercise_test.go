package main

import (
	"encoding/json"
	"testing"
)

func TestSafeMapRead(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	v, ok := SafeMapRead(m, "a")
	if !ok || v != 1 { t.Error() }
	v, ok = SafeMapRead(m, "z")
	if ok { t.Error() }
}

func TestPersonToJSON(t *testing.T) {
	p := Person{Name: Name{First: "John", Last: "Doe"}, Age: 30}
	b, err := PersonToJSON(p)
	if err != nil { t.Fatal(err) }
	var r map[string]interface{}
	json.Unmarshal(b, &r)
	if r["first_name"] != "John" || r["last_name"] != "Doe" || r["age"] != float64(30) {
		t.Errorf("bad json: %v", r)
	}
}
