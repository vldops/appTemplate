package main

import (
	"log"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	*zap.Logger
	error error
}

func (l *logger) makeOne() {
	cfg := zap.Config{
		Encoding:         "json",
		Development:      true,
		Level:            zap.NewAtomicLevelAt(0),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{

			MessageKey:       "message",
			ConsoleSeparator: " ",
			LevelKey:         "level",
			EncodeLevel:      zapcore.CapitalLevelEncoder,
			TimeKey:          "time",
			EncodeTime:       zapcore.RFC3339TimeEncoder,
			CallerKey:        "caller",
			EncodeCaller:     zapcore.ShortCallerEncoder,
		},
		InitialFields: map[string]interface{}{
			"appName": "app",
		},
	}

	l.Logger, l.error = cfg.Build()

	if l.error != nil {
		log.Fatal(l.error)
	}

}

func (l *logger) genLogs(r *reQuest) {
	l.Info("metricsMessage",
		zap.Int("statusCode", r.statusCode),
		zap.String("method", r.method),
		zap.String("path", r.path),
		zap.Float64("requestTime", r.duration),
		zap.String("remoteAddr", r.remoteAddr),
	)
}

func (l *logger) grubLogs() func(http.Handler) http.Handler {
	return func(serve http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			request := makeReQuest(w)
			serve.ServeHTTP(request, r)
			request.parseMetrics(r)
			l.genLogs(request)

		})
	}
}
