package internalhttp

import (
	"net/http"
	"strconv"
	"time"

	"golang.org/x/exp/slog"
)

type RequestInfo struct {
	ClientIP    string
	Date        string
	Method      string
	Path        string
	HTTPVersion string
	Latency     string
	UserAgent   string
}

func middlewareLogging(log *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		info := requestInformation(r, duration)

		log.Info("Request info",
			slog.String("client ip", info.ClientIP),
			slog.String("date", info.Date),
			slog.String("method", info.Method),
			slog.String("method path", info.Path),
			slog.String("HTTP version", info.HTTPVersion),
			slog.String("processing time", info.Latency),
			slog.String("user agent", info.UserAgent),
		)
	}
}

func requestInformation(r *http.Request, duration time.Duration) RequestInfo {
	clientIP := r.RemoteAddr
	date := time.Now().Format("02/Jan/2006:15:04:05 -0700")
	method := r.Method
	path := r.URL.Path
	HTTPVersion := r.Proto
	latency := strconv.Itoa(int(duration.Milliseconds())) + "ms"
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "empty"
	}
	return RequestInfo{
		ClientIP:    clientIP,
		Date:        date,
		Method:      method,
		Path:        path,
		HTTPVersion: HTTPVersion,
		Latency:     latency,
		UserAgent:   userAgent,
	}
}
