package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type prometheusMetrics struct {
	rCounter   *prometheus.CounterVec
	rHistogram *prometheus.HistogramVec
}

func (p *prometheusMetrics) makeOne() {

	p.rCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requestsTotal",
			Help: "requestsTotal",
		},
		[]string{"method", "path", "statusCode"},
	)

	p.rHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "requestsDuration",
			Help:    "requestsDuration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "statusCode"},
	)

	prometheus.MustRegister(p.rCounter)
	prometheus.MustRegister(p.rHistogram)

}

func (p *prometheusMetrics) grubMetrics() func(http.Handler) http.Handler {
	return func(serve http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			request := makeReQuest(w)
			serve.ServeHTTP(request, r)
			request.parseMetrics(r)

			p.rCounter.With(prometheus.Labels{
				"method":     request.method,
				"path":       request.path,
				"statusCode": request.statusCodePrometheus,
			}).Inc()

			p.rHistogram.With(prometheus.Labels{
				"method":     request.method,
				"path":       request.path,
				"statusCode": request.statusCodePrometheus,
			}).Observe(request.duration)

		})
	}
}
