# Hints for 18

## TODO 1
```go
if n < 2 { return false }
for i := 2; i*i <= n; i++ {
    if n%i == 0 { return false }
}
return true
```

## TODO 2
```go
runes := []rune(s)
for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
    runes[i], runes[j] = runes[j], runes[i]
}
return string(runes)
```

## TODO 3
```go
if n <= 1 { return n }
a, b := 0, 1
for i := 2; i <= n; i++ {
    a, b = b, a+b
}
return b
```
