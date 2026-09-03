# Hints for 16

## TODO 1
```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
```

## TODO 2
```go
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    next.ServeHTTP(w, r)
})
```

## TODO 3
```go
name := r.URL.Query().Get("name")
if name == "" {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]string{"error": "name required"})
    return
}
json.NewEncoder(w).Encode(map[string]string{"message": "Hello, " + name})
```
