package grpc

import (
	"context"
	"fmt"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"time"
)

const logPath = "./logging/logging.txt"

func loggingInterceptor(log logger.Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		log.Info("Request info GRPC",
			slog.String("method", info.FullMethod),
			slog.String("processing time", duration.String()),
		)

		logInFileString := fmt.Sprintf("GRPC: %s %s", info.FullMethod, duration)
		if err := log.WriteLogInFile(logPath, logInFileString); err != nil {
			log.Error(fmt.Sprintf("error wriging log in file with path %s: %s", logPath, err.Error()))
		}

		return resp, err
	}
}
