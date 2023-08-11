package logger

//go:generate mockgen -source=logger.go -destination=mock/mock.go logger

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

type MyLogger struct {
	log *slog.Logger
}

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	WriteLogInFile(path string, result string) error
}

func NewLogger(level string, representation string) *MyLogger {
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

	return &MyLogger{
		log: log,
	}
}

func (l *MyLogger) WriteLogInFile(path string, result string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(result + "\n"); err != nil {
		return err
	}
	return nil
}

func (l *MyLogger) Info(msg string, args ...any) {
	l.log.Info(msg, args...)
}

func (l *MyLogger) Error(msg string, args ...any) {
	l.log.Error(msg, args...)
}
