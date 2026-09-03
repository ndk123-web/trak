# Hints for 13

## TODO 1
```go
ch := make(chan int, 2)
ch <- 10
ch <- 20
return <-ch + <-ch
```

## TODO 2
```go
select {
case v := <-ch:
    return v
case <-time.After(50 * time.Millisecond):
    return "timeout"
}
```

## TODO 3
```go
jobs := make(chan int, len(inputs))
results := make(chan int, len(inputs))
for w := 0; w < 2; w++ {
    go func() {
        for j := range jobs {
            results <- j * 2
        }
    }()
}
for _, v := range inputs { jobs <- v }
close(jobs)
out := make([]int, 0, len(inputs))
for i := 0; i < len(inputs); i++ { out = append(out, <-results) }
return out
```
