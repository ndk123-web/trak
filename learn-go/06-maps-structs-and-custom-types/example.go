package main

import (
	"encoding/json"
	"fmt"
)

type BaseEntity struct {
	ID int `json:"id"`
}

type User struct {
	BaseEntity
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

func main() {
	m := make(map[string]int)
	m["score"] = 100
	if v, ok := m["score"]; ok {
		fmt.Println("Found:", v)
	}
	u := User{BaseEntity: BaseEntity{ID: 101}, Name: "Navnath"}
	data, _ := json.Marshal(u)
	fmt.Println("JSON:", string(data))
}
