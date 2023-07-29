package grpc

import (
	event_pb "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/grpc/pb/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

type ServerGRPCConfig struct {
	Host              string
	Port              string
	MaxConnectionIdle time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
	Timeout           time.Duration
}

type ServerGRPC struct {
	srv     *grpc.Server
	handler HandlerGRPC
}

func NewServerGRPC(handler *HandlerGRPC, cfg ServerGRPCConfig) *ServerGRPC {
	serverOptions := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()),
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
		handler: *handler,
	}
}

func (s *ServerGRPC) Start(cfg ServerGRPCConfig) error {
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
