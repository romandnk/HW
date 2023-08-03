package internalhttp

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
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

func loggerMiddleware(log logger.Logger, logPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		info := requestInformation(c.Request, duration)

		log.Info("Request info HTTP",
			slog.String("client ip", info.ClientIP),
			slog.String("date", info.Date),
			slog.String("method", info.Method),
			slog.String("method path", info.Path),
			slog.String("HTTP version", info.HTTPVersion),
			slog.String("status code", strconv.Itoa(c.Writer.Status())),
			slog.String("processing time", info.Latency),
			slog.String("user agent", info.UserAgent),
		)

		logInFileString := fmt.Sprintf("HTTP: %s %s %s %s %s %d %s %s",
			info.ClientIP,
			info.Date,
			info.Method,
			info.Path,
			info.HTTPVersion,
			c.Writer.Status(),
			info.Latency,
			info.UserAgent,
		)
		if err := log.WriteLogInFile(logPath, logInFileString); err != nil {
			log.Error(fmt.Sprintf("errors wriging log in file with path %s: %s", logPath, err.Error()))
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
