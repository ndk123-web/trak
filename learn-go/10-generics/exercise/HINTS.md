# Hints for 10

## TODO 1
```go
type Ordered interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64
}
func Max[T Ordered](a, b T) T {
    if a > b { return a }
    return b
}
```

## TODO 2
```go
out := make([]T, 0)
for _, v := range slice {
    if predicate(v) { out = append(out, v) }
}
return out
```

## TODO 3
```go
type Stack[T any] struct { items []T }
func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }
func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 { var z T; return z, false }
    v := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return v, true
}
```
