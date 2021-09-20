package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type setterHandler struct {
	set func(key string, value interface{})
}

func newSetterHandler(dic *diContainer) *setterHandler {
	return &setterHandler{
		set: dic.cache().Set,
	}
}

func newSetterHandlerDIProvider(dic *diContainer) func() (http.Handler, error) {
	var s *setterHandler
	var mu sync.Mutex
	return func() (http.Handler, error) {
		mu.Lock()
		defer mu.Unlock()
		if s == nil {
			s = newSetterHandler(dic)
		}
		return s, nil
	}
}

func configureSetHTTPRoute(r *mux.Route) *mux.Route {
	return r.Methods(http.MethodPost).Path("/set")
}

func (h *setterHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handle(req.Context(), w, req)
}

type setterHTTPHandlerRequestBody struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (h *setterHandler) handle(_ context.Context, w http.ResponseWriter, req *http.Request) {
	body, err := h.decodeRequestBody(req)
	if err != nil {
		onHTTPError(
			req.Context(),
			w,
			req,
			err,
			&httpErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
		return
	}
	// set key in prefix trie
	h.set(body.Key, body.Value)

	// set reverse key in prefix trie for suffix search without "data"
	h.set(reverse(body.Key), nil)

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("accepted"))
}

func (h *setterHandler) decodeRequestBody(req *http.Request) (*setterHTTPHandlerRequestBody, error) {
	var v *setterHTTPHandlerRequestBody
	err := json.NewDecoder(req.Body).Decode(&v)
	if err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return v, nil
}
