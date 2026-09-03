package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	HealthHandler(w, req)
	if w.Code != 200 { t.Errorf("code=%d", w.Code) }
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["status"] != "healthy" { t.Errorf("body=%v", body) }
}

func TestJSONMiddleware(t *testing.T) {
	handler := JSONMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}"))
	}))
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Header().Get("Content-Type") != "application/json" { t.Error() }
}

func TestGreetHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/greet?name=Alice", nil)
	w := httptest.NewRecorder()
	GreetHandler(w, req)
	if w.Code != 200 { t.Errorf("code=%d", w.Code) }
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != "Hello, Alice" { t.Errorf("body=%v", body) }

	req2 := httptest.NewRequest("GET", "/greet", nil)
	w2 := httptest.NewRecorder()
	GreetHandler(w2, req2)
	if w2.Code != 400 { t.Errorf("code=%d", w2.Code) }
}
