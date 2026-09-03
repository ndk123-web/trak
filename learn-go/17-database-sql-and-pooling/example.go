package main

import (
	"fmt"
	"time"
)

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func main() {
	cfg := PoolConfig{MaxOpenConns: 25, MaxIdleConns: 25, ConnMaxLifetime: 5 * time.Minute}
	fmt.Printf("Pool Config: %+v\n", cfg)
}
