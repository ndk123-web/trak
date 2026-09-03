# Hints for 08

## TODO 1
```go
type Rectangle struct { Width, Height float64 }
func (r Rectangle) Area() float64 { return r.Width * r.Height }
func (r *Rectangle) Scale(f float64) { r.Width *= f; r.Height *= f }
```

## TODO 2
```go
temp := *a; *a = *b; *b = temp
```

## TODO 3
```go
out := make([]*Rectangle, len(items))
for i := range items { out[i] = &items[i] }
return out
```
