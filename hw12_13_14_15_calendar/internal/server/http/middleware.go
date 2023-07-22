package internalhttp

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"golang.org/x/exp/slog"
)

const (
	EmptyStatusCode = "empty"
	logPath         = "./logging/logging.txt"
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

type statusCodeRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusCodeRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func middlewareLogging(log *logger.MyLogger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusCodeRecorder{ResponseWriter: w}

		start := time.Now()

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)

		info := requestInformation(r, duration)
		statusCode := processStatusCode(recorder.statusCode)

		log.Info("Request info",
			slog.String("client ip", info.ClientIP),
			slog.String("date", info.Date),
			slog.String("method", info.Method),
			slog.String("method path", info.Path),
			slog.String("HTTP version", info.HTTPVersion),
			slog.String("status code", statusCode),
			slog.String("processing time", info.Latency),
			slog.String("user agent", info.UserAgent),
		)

		logInFileString := fmt.Sprintf("%s %s %s %s %s %s %s %s",
			info.ClientIP,
			info.Date,
			info.Method,
			info.Path,
			info.HTTPVersion,
			statusCode,
			info.Latency,
			info.UserAgent,
		)
		if err := log.WriteLogInFile(logPath, logInFileString); err != nil {
			log.Error(fmt.Sprintf("error wriging log in file with path %s: %s", logPath, err.Error()))
		}
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

func processStatusCode(code int) string {
	if code == 0 {
		return EmptyStatusCode
	}
	return strconv.Itoa(code)
}
