package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type reQuest struct {
	http.ResponseWriter
	startAt              time.Time
	path                 string
	method               string
	remoteAddr           string
	duration             float64
	statusCode           int
	statusCodePrometheus string
}

func makeReQuest(w http.ResponseWriter) *reQuest {
	return &reQuest{
		ResponseWriter: w,
		startAt:        time.Now(),
	}
}

func (r *reQuest) parseMetrics(request *http.Request) {
	r.path = request.URL.EscapedPath()
	r.method = request.Method
	r.remoteAddr = strings.Split(request.RemoteAddr, ":")[0]
	r.duration = time.Since(r.startAt).Seconds()
	r.statusCodePrometheus = strconv.Itoa(r.statusCode)
}

// WriteHeader method
func (r *reQuest) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
