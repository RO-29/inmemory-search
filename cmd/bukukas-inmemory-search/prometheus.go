package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of total requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)
		statusCode := rw.statusCode
		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}

func initPromethus() {
	err := registerPrometheus()
	if err != nil {
		log.Fatal(err)
	}
	promRoute := mux.NewRouter()
	promRoute.Use(prometheusMiddleware)
	promRoute.Path("/prometheus").Handler(promhttp.Handler())
	fmt.Println("Serving requests on port :9000 for promethus")
	go func() {
		_ = http.ListenAndServe(":9000", promRoute)
	}()
}

func registerPrometheus() error {
	err := prometheus.Register(totalRequests)
	if err != nil {
		return errors.Wrap(err, "total requests")
	}
	err = prometheus.Register(responseStatus)
	if err != nil {
		return errors.Wrap(err, "response status")
	}
	err = prometheus.Register(httpDuration)
	if err != nil {
		return errors.Wrap(err, "http duration")
	}
	return nil
}
