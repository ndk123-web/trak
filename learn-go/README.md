# 🚀 Complete Go (Golang) End-to-End Learning Workspace

Welcome to your full Go mastery track materialized with **Trak**.

## 📚 Comprehensive Module Index

| Module | Title | Core Concepts |
| :--- | :--- | :--- |
| `00-setup-and-prerequisites` | **Setup & Environment** | Installation, GOPATH/GOROOT, VS Code tools |
| `01-runtime-and-escape-analysis` | **Go Runtime & Memory** | TCMalloc, Tri-Color GC, Stack vs Heap, escape analysis |
| `02-toolchain-and-workspaces` | **Toolchain & Modules** | `go mod`, cross-compilation, workspaces, race detector |
| `03-variables-and-zero-values` | **Variables & Types** | Zero values, short declarations, `iota` enums/bitmasks |
| `04-data-types-and-memory-headers` | **Memory Layout** | String/Slice headers, byte vs rune |
| `05-arrays-slices-and-growth` | **Slices Deep-Dive** | Growth algorithm, sub-slice leaks, `copy()` |
| `06-maps-structs-and-custom-types` | **Maps & Structs** | Hash maps, comma-ok, embedding, JSON tags |
| `07-control-flow-and-functions` | **Control Flow** | For range, defer LIFO, panic & recover |
| `08-pointers-and-semantics` | **Pointers & Semantics** | Pass-by-value, pointer vs value receivers |
| `09-methods-interfaces-duck-typing` | **Interfaces** | Duck typing, `iface` vs `eface`, nil interface trap |
| `10-generics` | **Generics (Go 1.18+)** | Type parameters, constraints, `any`, `comparable` |
| `11-error-handling-and-custom-types` | **Error Architecture** | Custom errors, `%w` wrapping, `errors.Is/As` |
| `12-concurrency-goroutines-gmp` | **Goroutines & GMP** | GMP model, cooperative scheduling, atomic ops |
| `13-channels-select-and-sync` | **Channels & Sync** | Buffered/unbuffered, select timeout, sync primitives |
| `14-context-cancellation-timeouts` | **Context Tree** | `WithCancel`, `WithTimeout`, `WithValue`, cancellation |
| `15-standard-library-powerhouses` | **Standard Library** | `io.Reader/Writer`, `encoding/json`, time formatting |
| `16-production-web-services` | **Web Architecture** | `net/http`, middleware, structured JSON APIs, timeouts |
| `17-database-sql-and-pooling` | **Database & Pooling** | `database/sql`, connection pool tuning, null types |
| `18-testing-benchmarks-profiling` | **Testing & Profiling** | Table-driven tests, subtests, benchmarks, race detection |
| `19-interview-questions-and-drills` | **Interview Master Drill** | Top 10 high-yield interview questions with runnable checks |

---

## 🛠️ How to Navigate & Learn

Each module has this structure:
```
XX-module-name/
├── README.md              # Deep learning material + interview questions
├── example.go             # Reference code (runnable demo)
├── exercise/
│   ├── starter.go         # Incomplete code with TODOs — YOU fill this
│   └── HINTS.md           # Progressive hints (light → detailed)
└── verify/
    ├── exercise_test.go   # Auto-verification: run `go test ./...`
    └── solution.go        # Complete answer (check AFTER attempting)
```

### Learning Flow
1. **Read** `README.md` to understand concepts and traps
2. **Run** `example.go` to see the concept in action
3. **Code** `exercise/starter.go` — fill in the TODOs
4. **Verify** with `go test ./verify/...` — green = understanding
5. **Check** `verify/solution.go` only if stuck

### Quick Commands
```bash
# Verify your solution
cd XX-module-name
go test ./verify/...

# Run with race detector
go test -race ./verify/...

# Run example
go run example.go
```

---

## 🎯 Verification Strategy

Every module includes `verify/exercise_test.go` that you can run with:
```bash
go test ./verify/...
```

- **PASS** = Your code is correct. Move to the next module.
- **FAIL** = Re-read the README, check HINTS.md, debug your starter.go.
- **RACE** = Run `go test -race ./verify/...` to catch data races.

---

## 🏆 Completion Criteria

A module is "complete" when:
1. ✅ You can explain the core concept without looking at notes
2. ✅ `go test ./verify/...` passes all tests
3. ✅ `go test -race ./verify/...` passes (no data races)
4. ✅ You can answer the interview questions at the bottom of README.md

Happy learning! 🎉
