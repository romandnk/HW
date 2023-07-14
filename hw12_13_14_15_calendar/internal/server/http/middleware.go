package internalhttp

import (
	"golang.org/x/exp/slog"
	"net/http"
	"time"
)

func middlewareLogging(log *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		log.Info("Request info",
			slog.String("method", r.Method),
			slog.String("method path", r.URL.Path),
			slog.Duration("processing time", duration),
		)
	}
}
