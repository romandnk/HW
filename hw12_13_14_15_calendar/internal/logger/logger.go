package logger

import (
	"os"

	"golang.org/x/exp/slog"
)

const (
	infoLevel  = "INFO"
	debugLevel = "DEBUG"
	errorLevel = "ERROR"
	warnLevel  = "WARN"
	jsonLogger = "JSON"
	textLogger = "TEXT"
)

func NewLogger(level string, representation string) *slog.Logger {
	var log *slog.Logger

	logOptions := slog.HandlerOptions{}

	switch level {
	case infoLevel:
		logOptions.Level = slog.LevelInfo
	case debugLevel:
		logOptions.Level = slog.LevelDebug
	case errorLevel:
		logOptions.Level = slog.LevelError
	case warnLevel:
		logOptions.Level = slog.LevelWarn
	}

	switch representation {
	case jsonLogger:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &logOptions))
	case textLogger:
		log = slog.New(slog.NewTextHandler(os.Stdout, &logOptions))
	}

	return log
}
