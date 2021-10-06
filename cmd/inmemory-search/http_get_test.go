package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetHandlerFound(t *testing.T) {
	r := mux.NewRouter()
	var getCache bool
	h := &getterHandler{
		get: func(key string) (value interface{}, found bool) {
			getCache = true
			return "1", true
		},
	}
	configureGetHTTPRoute(r.NewRoute()).Handler(h)
	req := httptest.NewRequest(http.MethodGet, "/get/xyz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Log(w.Body)
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusOK)
	}
	if !getCache {
		t.Fatalf("unexpected get cache: got %#v, want %#v", getCache, true)
	}
}

func TestGetHandlerNotFound(t *testing.T) {
	r := mux.NewRouter()
	var getCache bool
	h := &getterHandler{
		get: func(key string) (value interface{}, found bool) {
			getCache = true
			return "1", false
		},
	}
	configureGetHTTPRoute(r.NewRoute()).Handler(h)
	req := httptest.NewRequest(http.MethodGet, "/get/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	w.Flush()

	if w.Code != http.StatusNotFound {
		t.Log(w.Body)
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusNotFound)
	}
	if !getCache {
		t.Fatalf("unexpected get cache: got %#v, want %#v", getCache, true)
	}
}
