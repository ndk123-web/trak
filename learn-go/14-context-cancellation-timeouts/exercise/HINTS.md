# Hints for 14

## TODO 1
```go
ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
defer cancel()
select {
case <-time.After(50 * time.Millisecond):
    return false
case <-ctx.Done():
    return true
}
```

## TODO 2
```go
parent, cancel := context.WithCancel(context.Background())
child, _ := context.WithTimeout(parent, time.Minute)
cancel()
select {
case <-child.Done():
    return true
case <-time.After(100 * time.Millisecond):
    return false
}
```

## TODO 3
```go
ctx = context.WithValue(ctx, requestIDKey{}, "req-123")
if id, ok := ctx.Value(requestIDKey{}).(string); ok {
    return id
}
return ""
```
