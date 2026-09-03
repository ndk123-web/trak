# Hints for 17

## TODO 1
```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
return db
```

## TODO 2
```go
if ns.Valid {
    return ns.String
}
return "N/A"
```

## TODO 3
```go
return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)
```
