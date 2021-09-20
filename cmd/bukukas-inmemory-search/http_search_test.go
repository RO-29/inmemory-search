package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestSearchHandlerFoundPrefix(t *testing.T) {
	r := mux.NewRouter()
	var searchP bool
	h := &searchHandler{
		search: func(keyPattern string, total int) []string {
			searchP = true
			return []string{"1", "2"}
		},
	}
	configureSearchHTTPRoute(r.NewRoute()).Handler(h)
	req := httptest.NewRequest(http.MethodGet, "/search?prefix=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Log(w.Body)
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusOK)
	}
	if !searchP {
		t.Fatalf("unexpected get cache: got %#v, want %#v", searchP, true)
	}
}

func TestSearchHandlerFoundSuffix(t *testing.T) {
	r := mux.NewRouter()
	var searchP bool
	h := &searchHandler{
		search: func(keyPattern string, total int) []string {
			searchP = true
			return []string{"2", "4"}
		},
	}
	configureSearchHTTPRoute(r.NewRoute()).Handler(h)
	req := httptest.NewRequest(http.MethodGet, "/search?suffix=5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Log(w.Body)
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusOK)
	}
	if !searchP {
		t.Fatalf("unexpected get cache: got %#v, want %#v", searchP, true)
	}
}

func TestSearchHandlerBadSearch(t *testing.T) {
	r := mux.NewRouter()
	var searchP bool
	h := &searchHandler{
		search: func(keyPattern string, total int) []string {
			searchP = true
			return []string{"2", "4"}
		},
	}
	configureSearchHTTPRoute(r.NewRoute()).Handler(h)
	req := httptest.NewRequest(http.MethodGet, "/search?rand=5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusBadRequest {
		t.Log(w.Body)
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusOK)
	}
	if searchP {
		t.Fatalf("unexpected get cache: got %#v, want %#v", searchP, false)
	}
}
