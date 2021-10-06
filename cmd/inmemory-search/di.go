package main

import (
	"inmemory-search/inmemorycache"
	"net/http"

	"github.com/gorilla/mux"
)

type diContainer struct {
	flags *flags

	httpHandler  func() (http.Handler, error)
	httpRouter   func() (*mux.Router, error)
	httpHandlers *httpHandlers
	cache        func() inmemorycache.Cache
}

func newDIContainer(flg *flags) *diContainer {
	dic := &diContainer{
		flags: flg,
	}
	dic.httpHandlers = newHTTPHandlers(dic)
	dic.httpRouter = newHTTPRouterDIProvider(dic)
	dic.httpHandler = newHTTPHandlerDIProvider(dic)
	dic.cache = inmemorycache.NewDIProvider()
	return dic
}
