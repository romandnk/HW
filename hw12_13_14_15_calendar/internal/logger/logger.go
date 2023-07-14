package logger

import (
	"golang.org/x/exp/slog"
	"os"
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

	if representation == jsonLogger {
		switch level {
		case infoLevel:
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		case debugLevel:
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		case errorLevel:
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		case warnLevel:
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
		}
	} else if representation == textLogger {
		switch level {
		case infoLevel:
			log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		case debugLevel:
			log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		case errorLevel:
			log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		case warnLevel:
			log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
		}
	}

	return log
}
