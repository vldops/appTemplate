package main

import (
	"encoding/json"
	"net/http"
	"time"

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
	gorillaMux.HandleFunc("/one", sleepOne)
	gorillaMux.HandleFunc("/two", sleepTwo)
	gorillaMux.HandleFunc("/three", sleepThree)
	gorillaMux.HandleFunc("/300ms", sleep300ms)
	gorillaMux.Handle("/metrics", promhttp.Handler())
	gorillaMux.Use(log.grubLogs(), prometheus.grubMetrics())
	log.Fatal("serverErrorMsg", zap.Error(http.ListenAndServe(":8080", gorillaMux)))
}

func goOk(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"/": true})

}

func defaultRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{r.URL.EscapedPath(): "goFuckYourself"})
}

func sleepOne(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{r.URL.EscapedPath(): "ok"})
}

func sleepTwo(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{r.URL.EscapedPath(): "ok"})
}

func sleepThree(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{r.URL.EscapedPath(): "ok"})
}

func sleep300ms(w http.ResponseWriter, r *http.Request) {
	time.Sleep(300 * time.Millisecond)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{r.URL.EscapedPath(): "ok"})
}
