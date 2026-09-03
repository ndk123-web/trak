# 16 — Production HTTP Web Services

## 🎯 Learning Objectives
- Build structured JSON APIs with `net/http`
- Implement middleware chains
- Handle graceful shutdown with `http.Server`
- Set appropriate timeouts (ReadTimeout, WriteTimeout, IdleTimeout)

---

## 🧠 Core Concepts

### http.ServeMux (Go 1.22+)
```go
mux := http.NewServeMux()
mux.HandleFunc("GET /health", handler) // method-specific patterns
mux.HandleFunc("/api/", apiHandler)    // subtree matching
```

### Middleware Pattern
```go
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL, time.Since(start))
    })
}
```

### Graceful Shutdown
```go
srv := &http.Server{Addr: ":8080", Handler: mux}
go srv.ListenAndServe()
// wait for interrupt signal
srv.Shutdown(ctx)
```

### Timeouts
- `ReadTimeout`: Max time to read request body
- `WriteTimeout`: Max time to write response
- `IdleTimeout`: Max time between requests on keep-alive connection

---

## ⚠️ Common Traps
- **No timeouts**: Default server has NO timeouts — vulnerable to slowloris attacks
- **Not checking request method**: `HandleFunc` matches all methods unless specified
- **Writing header after body**: `w.WriteHeader()` after `w.Write()` is ignored

---

## 📝 Interview Questions
Q: What's the difference between `http.HandleFunc` and `mux.HandleFunc`?
A: `http.HandleFunc` registers on the default mux (`http.DefaultServeMux`). `mux.HandleFunc` registers on a specific mux instance.

Q: Why is `WriteTimeout` important?
A: Prevents slow clients from holding connections open indefinitely. Without it, a malicious client can exhaust file descriptors.
