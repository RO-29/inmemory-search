package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type pingHandler struct { // TODO Have fun here
}

func newPingHandler() *pingHandler {
	return &pingHandler{}
}

func newPingHandlerDIProvider() func() (http.Handler, error) {
	var s *pingHandler
	var mu sync.Mutex
	return func() (http.Handler, error) {
		mu.Lock()
		defer mu.Unlock()
		if s == nil {
			s = newPingHandler()
		}
		return s, nil
	}
}

func configurePingHTTPRoute(r *mux.Route) *mux.Route {
	return r.Methods(http.MethodGet).Path("/_ping")
}

func (h *pingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handle(req.Context(), w, req)
}

func (h *pingHandler) handle(_ context.Context, w http.ResponseWriter, _ *http.Request) {
	// TODO ideas for health check?
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("All well"))
}
