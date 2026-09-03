# Hints for 15

## TODO 1
```go
return io.Copy(w, r)
```

## TODO 2
```go
buf := make([]byte, n)
_, err := io.ReadFull(r, buf)
return string(buf), err
```

## TODO 3
```go
return io.MultiReader(r1, r2)
```
