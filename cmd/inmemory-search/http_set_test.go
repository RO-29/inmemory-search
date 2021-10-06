package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestSetHandlerSet(t *testing.T) {
	r := mux.NewRouter()
	var setCache bool
	h := &setterHandler{
		set: func(key string, value interface{}) {
			setCache = true
		},
	}
	body := `{
		"key": "abc",
		"value": "1"
	}`
	configureSetHTTPRoute(r.NewRoute()).Handler(h)
	req := httptest.NewRequest(http.MethodPost, "/set", strings.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusCreated {
		t.Log(w.Body)
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusCreated)
	}
	if !setCache {
		t.Fatalf("unexpected set cache: got %#v, want %#v", setCache, true)
	}
}

func TestSetHandlerSetErr(t *testing.T) {
	r := mux.NewRouter()
	var setCache bool
	h := &setterHandler{
		set: func(key string, value interface{}) {
			setCache = true
		},
	}
	body := `{
		"keys": "abc",
		"value": "1"
	}`
	configureSetHTTPRoute(r.NewRoute()).Handler(h)
	req := httptest.NewRequest(http.MethodPost, "/set", strings.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusBadRequest {
		t.Log(w.Body)
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusBadRequest)
	}
	if setCache {
		t.Fatalf("unexpected set cache: got %#v, want %#v", setCache, false)
	}
}
