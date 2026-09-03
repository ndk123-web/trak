package main

import (
	"database/sql"
	"strings"
	"testing"
	"time"
)

func TestConfigurePool(t *testing.T) {
	// We can't easily test internal state without a real DB,
	// but we can verify it doesn't panic.
	db, _ := sql.Open("mysql", "")
	_ = ConfigurePool(db)
}

func TestScanNullableString(t *testing.T) {
	ns := sql.NullString{String: "hello", Valid: true}
	if ScanNullableString(ns) != "hello" { t.Error() }
	ns2 := sql.NullString{Valid: false}
	if ScanNullableString(ns2) != "N/A" { t.Error() }
}

func TestBuildDSN(t *testing.T) {
	dsn := BuildDSN("root", "pass", "localhost", "3306", "mydb")
	if !strings.Contains(dsn, "root:pass") { t.Error() }
	if !strings.Contains(dsn, "@tcp(localhost:3306)") { t.Error() }
	if !strings.Contains(dsn, "/mydb?parseTime=true") { t.Error() }
}
