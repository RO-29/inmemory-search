package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type searchHandler struct {
	search func(keyPattern string, total int) []string
}

func newSearchHandler(dic *diContainer) *searchHandler {
	return &searchHandler{
		search: dic.cache().PrefixMatch,
	}
}

func newSearchHandlerDIProvider(dic *diContainer) func() (http.Handler, error) {
	var s *searchHandler
	var mu sync.Mutex
	return func() (http.Handler, error) {
		mu.Lock()
		defer mu.Unlock()
		if s == nil {
			s = newSearchHandler(dic)
		}
		return s, nil
	}
}

func configureSearchHTTPRoute(r *mux.Route) *mux.Route {
	return r.Methods(http.MethodGet).Path("/search")
}

func (h *searchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handle(req.Context(), w, req)
}

type searchHTTPHandlerResponseBody struct {
	Keys []string `json:"keys"`
}

func (h *searchHandler) handle(_ context.Context, w http.ResponseWriter, req *http.Request) {
	var keys []string
	var keyPattern string
	if req.URL.Query().Get("prefix") != "" { //nolint: gocritic // not useful suggestion
		keyPattern = req.URL.Query().Get("prefix")
		keys = h.handleQueryPattern(keyPattern, false)
	} else if req.URL.Query().Get("suffix") != "" {
		keyPattern = req.URL.Query().Get("suffix")
		keys = h.handleQueryPattern(keyPattern, true)
	} else {
		onHTTPError(
			req.Context(),
			w,
			req,
			errors.New("invalid search"),
			&httpErrorResponse{
				Code:    http.StatusBadRequest,
				Message: errors.New("invalid search").Error(),
			})
	}
	if len(keys) == 0 {
		onHTTPError(
			req.Context(),
			w,
			req,
			errors.New(fmt.Sprintf("key pattern '%v' not found!", keyPattern)),
			&httpErrorResponse{
				Code:    http.StatusNotFound,
				Message: errors.New(fmt.Sprintf("key pattern '%v' not found!", keyPattern)).Error(),
			})
		return
	}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	_ = enc.Encode(&searchHTTPHandlerResponseBody{
		Keys: keys,
	})
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func (h *searchHandler) handleQueryPattern(keyPattern string, isReverse bool) []string {
	if !isReverse {
		return h.search(keyPattern, -1)
	}
	// reverse the key to search for suffix
	keysR := h.search(reverse(keyPattern), -1)
	// reverse the result
	keys := make([]string, len(keysR))
	for i, v := range keysR {
		keys[i] = reverse(v)
	}
	return keys
}
