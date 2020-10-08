package main

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	var log logger
	log.makeOne()

	var prometheus prometheusMetrics
	prometheus.makeOne()

	gorillaMux := mux.NewRouter()
	gorillaMux.HandleFunc("/ok", goOk)
	gorillaMux.Handle("/metrics", promhttp.Handler())
	gorillaMux.Use(log.grubLogs(), prometheus.grubMetrics())
	log.Fatal("serverErrorMsg", zap.Error(http.ListenAndServe(":8080", gorillaMux)))
}

func goOk(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"/": true})

}
