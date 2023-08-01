package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

const logPath = "./logging/logging.txt"

func loggingInterceptor(log logger.Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) { //nolint:lll
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		logErr := err
		if logErr == nil {
			logErr = errors.New("empty")
		}

		log.Info("Request info GRPC",
			slog.String("method", info.FullMethod),
			slog.String("processing time", duration.String()),
			slog.String("errors", logErr.Error()),
		)

		logInFileString := fmt.Sprintf("GRPC: %s %s %s", info.FullMethod, duration, logErr.Error())
		if err := log.WriteLogInFile(logPath, logInFileString); err != nil {
			log.Error(fmt.Sprintf("errors wriging log in file with path %s: %s", logPath, err.Error()))
		}

		return resp, err
	}
}
