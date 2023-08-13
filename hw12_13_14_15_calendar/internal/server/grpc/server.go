package grpc

import (
	"net"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/cmd/config"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	event_pb "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/grpc/pb/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ServerGRPC struct {
	srv     *grpc.Server
	handler *HandlerGRPC
}

func NewServerGRPC(handler *HandlerGRPC, log logger.Logger, cfg config.ServerGRPCConfig, logPath string) *ServerGRPC {
	serverOptions := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(loggingInterceptor(log, logPath)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: cfg.MaxConnectionIdle,
			MaxConnectionAge:  cfg.MaxConnectionAge,
			Time:              cfg.Time,
			Timeout:           cfg.Timeout,
		}),
	}

	srv := grpc.NewServer(serverOptions...)

	return &ServerGRPC{
		srv:     srv,
		handler: handler,
	}
}

func (s *ServerGRPC) Start(cfg config.ServerGRPCConfig) error {
	lsn, err := net.Listen("tcp", net.JoinHostPort(cfg.Host, cfg.Port))
	if err != nil {
		return err
	}

	event_pb.RegisterEventServiceServer(s.srv, s.handler)

	return s.srv.Serve(lsn)
}

func (s *ServerGRPC) Stop() {
	s.srv.GracefulStop()
}
