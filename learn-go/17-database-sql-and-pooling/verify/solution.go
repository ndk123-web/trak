package main

import (
	"database/sql"
	"fmt"
	"time"
)

func ConfigurePool(db *sql.DB) *sql.DB {
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db
}

func ScanNullableString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return "N/A"
}

func BuildDSN(user, password, host, port, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)
}

func main() {}
