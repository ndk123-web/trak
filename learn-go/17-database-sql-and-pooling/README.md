# 17 — Database Access & Connection Pooling

## 🎯 Learning Objectives
- Use `database/sql` for queries and transactions
- Tune connection pool settings (`SetMaxOpenConns`, etc.)
- Handle `sql.NullString`, `sql.NullInt64` for nullable columns
- Understand prepared statements and their reuse

---

## 🧠 Core Concepts

### Connection Pool Settings
```go
db.SetMaxOpenConns(25)       // max simultaneous connections
db.SetMaxIdleConns(25)       // max idle connections in pool
db.SetConnMaxLifetime(5 * time.Minute) // max age before closing
```

### Null Types
```go
var name sql.NullString
if name.Valid { fmt.Println(name.String) }
```

### Transactions
```go
tx, err := db.Begin()
// ... queries ...
tx.Commit() // or tx.Rollback()
```

---

## ⚠️ Common Traps
- **Not closing rows**: `defer rows.Close()`
- **Using `db.Query` for INSERT**: Use `db.Exec` instead
- **Scanning into wrong types**: Causes `sql: Scan error`
- **Not handling NULL**: Use `sql.Null*` types or `COALESCE` in SQL

---

## 📝 Interview Questions
Q: What's the difference between `SetMaxOpenConns` and `SetMaxIdleConns`?
A: `MaxOpenConns` limits total active connections. `MaxIdleConns` limits how many of those can be idle in the pool. If `MaxIdleConns < MaxOpenConns`, connections above the idle limit are closed when returned to the pool.

Q: Should you use `db.Query` or `db.Exec` for INSERT?
A: `Exec`. `Query` is for SELECTs that return rows.
