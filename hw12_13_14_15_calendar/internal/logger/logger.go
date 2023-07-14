package logger

import (
	"golang.org/x/exp/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(level string) *Logger {
	var log *slog.Logger

	switch level {
	case "INFO":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case "DEBUG":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "ERROR":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	case "WARN":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	}

	return &Logger{
		Logger: log,
	}
}
