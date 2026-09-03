# 06 — Maps, Structs, Embedding & JSON Tags

## 🎯 Learning Objectives
- Use maps safely (NOT concurrency-safe!)
- Use comma-ok idiom for map lookups
- Use struct embedding for composition
- Marshal/unmarshal structs with JSON tags

---

## 🧠 Core Concepts

### Map Internals
- Hash tables with buckets
- **NOT safe for concurrent use** — concurrent read+write panics
- Zero value is `nil` — must `make()` before writing

### Struct Embedding
```go
type Animal struct { Name string }
type Dog struct {
    Animal      // embedded — promotes fields/methods
    Breed string
}
```

### JSON Tags
- `json:"name,omitempty"` — skips zero values
- `json:"-"` — ignores field entirely

---

## ⚠️ Common Traps
- Map iteration order is random
- Nil map read is safe, write panics
- Embedding promotes methods; named fields don't

---

## 📝 Interview Questions
Q: What happens writing to a map from two goroutines?
A: Runtime panic: `concurrent map writes`. Use `sync.RWMutex` or `sync.Map`.

Q: What's the difference between `json:"-"` and omitting the tag?
A: `json:"-"` explicitly excludes the field. No tag includes it with the field name.
