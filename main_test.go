package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_PingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected code 200 got %d", w.Code)
	}

	if w.Body.String() != "pong" {
		t.Errorf("Expected body to be 'pong' got %s", w.Body.String())
	}
}
