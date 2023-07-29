package grpc

import (
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	event_pb "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/grpc/pb/event"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/service"
)

type HandlerGRPC struct {
	event_pb.UnimplementedEventServiceServer
	service service.Services
	logger  logger.Logger
}

func NewHandlerGRPC(services service.Services, logger logger.Logger) *HandlerGRPC {
	return &HandlerGRPC{
		UnimplementedEventServiceServer: event_pb.UnimplementedEventServiceServer{},
		service:                         services,
		logger:                          logger,
	}
}
