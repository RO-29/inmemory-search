package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func runHTTPServer(_ context.Context, dic *diContainer, addr string) error {
	h, err := dic.httpHandler()
	if err != nil {
		return errors.Wrap(err, "get http handler")
	}
	log.Printf("Start HTTP server on %s", addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: h,
	}
	err = srv.ListenAndServe()
	if err != nil {
		return errors.Wrap(err, "listen and serve")
	}
	log.Println("Stopped HTTP server")
	return nil
}

func newHTTPHandler(dic *diContainer) (http.Handler, error) {
	r, err := dic.httpRouter()
	if err != nil {
		return nil, errors.Wrap(err, "router")
	}
	return r, nil
}

func newHTTPHandlerDIProvider(dic *diContainer) func() (http.Handler, error) {
	var v http.Handler
	var mu sync.Mutex
	return func() (_ http.Handler, err error) {
		mu.Lock()
		defer mu.Unlock()
		if v == nil {
			v, err = newHTTPHandler(dic)
		}
		return v, err
	}
}

func newHTTPRouter(dic *diContainer) (*mux.Router, error) {
	r := mux.NewRouter()
	err := registerHTTPHandlers(r, dic.httpHandlers)
	if err != nil {
		return nil, errors.Wrap(err, "register")
	}
	return r, nil
}

func newHTTPRouterDIProvider(dic *diContainer) func() (*mux.Router, error) {
	var r *mux.Router
	var mu sync.Mutex
	return func() (*mux.Router, error) {
		mu.Lock()
		defer mu.Unlock()
		var err error
		if r == nil {
			r, err = newHTTPRouter(dic)
		}
		return r, err
	}
}

func registerHTTPHandlers(r *mux.Router, hs *httpHandlers) error {
	for _, v := range []struct {
		name      string
		configure func(*mux.Route) *mux.Route
		handler   func() (http.Handler, error)
	}{
		{
			name:      "set",
			configure: configureSetHTTPRoute,
			handler:   hs.setterHandler,
		},
		{
			name: "get",
			// configure: configureRealtimeStatsHTTPRoute,
			// handler:   hs.RealtimeStatsHandler,
		},
		{
			name: "search",
			// configure: configureRealtimeStatsHTTPRoute,
			// handler:   hs.RealtimeStatsHandler,
		},
	} {
		h, err := v.handler()
		if err != nil {
			return errors.Wrap(err, v.name)
		}
		v.configure(r.NewRoute()).Handler(h)
	}
	return nil
}

type httpHandlers struct {
	setterHandler func() (http.Handler, error)
}

func newHTTPHandlers(dic *diContainer) *httpHandlers {
	return &httpHandlers{
		setterHandler: newSetterHandlerDIProvider(dic),
	}
}

func onHTTPError(_ context.Context, w http.ResponseWriter, _ *http.Request, err error, resp *httpErrorResponse) {
	hd := w.Header()
	hd.Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	_ = enc.Encode(resp)
	_, _ = w.Write(buf.Bytes())
	log.Println("err: ", err)
}

type httpErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
