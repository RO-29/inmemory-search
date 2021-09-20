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

type getterHandler struct {
	get func(key string) (value interface{}, found bool)
}

func newGetterHandler(dic *diContainer) *getterHandler {
	return &getterHandler{
		get: dic.cache().Get,
	}
}

func newGetterHandlerDIProvider(dic *diContainer) func() (http.Handler, error) {
	var s *getterHandler
	var mu sync.Mutex
	return func() (http.Handler, error) {
		mu.Lock()
		defer mu.Unlock()
		if s == nil {
			s = newGetterHandler(dic)
		}
		return s, nil
	}
}

func configureGetHTTPRoute(r *mux.Route) *mux.Route {
	return r.Methods(http.MethodGet).Path("/get/{key}")
}

func (h *getterHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handle(req.Context(), w, req)
}

type getterHTTPHandlerResponseBody struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (h *getterHandler) handle(_ context.Context, w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	key := vars["key"]
	if key == "" {
		onHTTPError(
			req.Context(),
			w,
			req,
			errors.New("empty key"),
			&httpErrorResponse{
				Code:    http.StatusBadRequest,
				Message: errors.New("empty key").Error(),
			})
		return
	}
	// get value from prefix trie
	value, ok := h.get(key)
	if !ok {
		onHTTPError(
			req.Context(),
			w,
			req,
			errors.New(fmt.Sprintf("key '%v' not found!", key)),
			&httpErrorResponse{
				Code:    http.StatusNotFound,
				Message: errors.New(fmt.Sprintf("key '%v' not found!", key)).Error(),
			})
		return
	}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	_ = enc.Encode(&getterHTTPHandlerResponseBody{
		Key:   key,
		Value: value,
	})
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
