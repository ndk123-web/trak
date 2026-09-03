# 15 — Standard Library: io, json, os, time

## 🎯 Learning Objectives
- Use `io.Reader`/`io.Writer` interfaces for streaming
- Marshal/unmarshal JSON efficiently
- Handle OS signals gracefully
- Parse and format time correctly

---

## 🧠 Core Concepts

### io.Reader / io.Writer
The most important interfaces in Go:
```go
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
```
Everything implements these: files, network connections, buffers, pipes.

### JSON Streaming
```go
enc := json.NewEncoder(w)
enc.Encode(obj)

dec := json.NewDecoder(r)
dec.Decode(&obj)
```

### Time Formatting
Go uses **reference time** for formatting:
```go
t.Format("2006-01-02 15:04:05") // YYYY-MM-DD HH:MM:SS
```
The reference time is: `Mon Jan 2 15:04:05 MST 2006` (1-2-3-4-5-6-7-8 pattern).

---

## ⚠️ Common Traps
- **JSON numbers decode to float64 by default**: Use `json.Decoder.UseNumber()` or custom types
- **Time zones**: `time.Now()` uses local time. `time.UTC()` converts.
- **io.EOF is not an error**: It signals normal end of stream

---

## 📝 Interview Questions
Q: What's the difference between `json.Marshal` and `json.NewEncoder`?
A: `Marshal` returns a `[]byte`. `Encoder` streams directly to an `io.Writer` without allocating the full byte slice.

Q: Why does Go use `2006-01-02 15:04:05` as the format reference?
A: It's a mnemonic: 1-2-3-4-5-6-7-8 (January=1, 2nd day, 3pm=15:00, 4 minutes, 5 seconds, 2006 year, -07:00 timezone offset).
