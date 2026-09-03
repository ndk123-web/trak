# Hints for 12

## TODO 1
```go
var counter int64
var wg sync.WaitGroup
for i := 0; i < n; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        atomic.AddInt64(&counter, 1)
    }()
}
wg.Wait()
return counter
```

## TODO 2
```go
var mu sync.Mutex
var result []int
var wg sync.WaitGroup
for i := 1; i <= 3; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        mu.Lock()
        result = append(result, id)
        mu.Unlock()
    }(i)
}
wg.Wait()
return result
```

## TODO 3
- `runtime.GOMAXPROCS(0)`
