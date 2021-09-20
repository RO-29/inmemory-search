package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var version = "dev"

const (
	appName = "bukukas-inmemory-search"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(errors.Wrap(err, "run"))
	}
}

func run() error {
	log.Printf("running version: %#v, app: %#v", version, appName)
	flg := getFlags()
	ctx := context.Background()
	initPromethus()
	dic := newDIContainer(flg)
	err := runHTTPServer(ctx, dic, flg.http)
	if err != nil {
		return errors.Wrap(err, "HTTP server")
	}
	return nil
}

func initPromethus() {
	promRoute := mux.NewRouter()
	promRoute.Path("/prometheus").Handler(promhttp.Handler())
	fmt.Println("Serving requests on port 9000 for promethus")
	go func() {
		_ = http.ListenAndServe(":9000", promRoute)
	}()
}
