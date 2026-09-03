# Hints for 07

## TODO 1
```go
fmt.Println("start")
defer fmt.Println("end")
fmt.Println("middle")
```

## TODO 2
```go
defer func() {
    if r := recover(); r != nil {
        err = fmt.Errorf("index out of range")
    }
}()
return slice[idx], nil
```

## TODO 3
```go
sum := 0
for i := 1; i <= 3; i++ {
    defer func(n int) { sum += n }(i)
}
return sum
```
