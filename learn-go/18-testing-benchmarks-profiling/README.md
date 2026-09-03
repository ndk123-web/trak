# 18 — Testing, Benchmarks, Table-Driven Tests & Race Detection

## 🎯 Learning Objectives
- Write table-driven tests with `t.Run()`
- Write benchmarks with `testing.B`
- Use subtests for better failure reporting
- Run the race detector (`go test -race`)

---

## 🧠 Core Concepts

### Table-Driven Tests
```go
func TestAdd(t *testing.T) {
    tests := []struct{
        name string
        a, b, want int
    }{
        {"positive", 2, 3, 5},
        {"negative", -2, -3, -5},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want { t.Errorf(...) }
        })
    }
}
```

### Benchmarks
```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
```

### Race Detector
```bash
go test -race ./...
```
Adds ~5-10x overhead. Detects data races at runtime.

---

## ⚠️ Common Traps
- **Not using `t.Parallel()`**: Tests run sequentially by default
- **Benchmark modifying global state**: Causes incorrect measurements
- **Forgetting `b.ResetTimer()`**: Includes setup time in benchmark

---

## 📝 Interview Questions
Q: What's the difference between `t.Error` and `t.Fatal`?
A: `t.Error` logs and continues. `t.Fatal` logs and stops the current test immediately.

Q: How do you benchmark memory allocations?
A: `b.ReportAllocs()` or run with `go test -benchmem`.
