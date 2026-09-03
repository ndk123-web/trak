# Hints for 06

## TODO 1
- `val, ok := m[key]; return val, ok`

## TODO 2
```go
type Name struct {
    First string `json:"first_name"`
    Last  string `json:"last_name"`
}
type Person struct {
    Name
    Age int `json:"age"`
}
```

## TODO 3
- `return json.Marshal(p)`
