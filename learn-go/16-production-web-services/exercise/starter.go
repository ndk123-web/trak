package main

import (
	"encoding/json"
	"net/http"
)

// TODO 1: Create a health check handler that writes JSON {"status":"healthy"}.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// FILL HERE
}

// TODO 2: Create a middleware that sets Content-Type to application/json.
func JSONMiddleware(next http.Handler) http.Handler {
	// FILL HERE
	return nil
}

// TODO 3: Create a handler that reads a "name" query parameter and returns
// JSON {"message":"Hello, <name>"}. If name is missing, return 400.
func GreetHandler(w http.ResponseWriter, r *http.Request) {
	// FILL HERE
}

func main() {}
