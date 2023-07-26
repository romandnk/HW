package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"
)

type ServerHTTP struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server struct {
	srv *http.Server
}

func NewServer(cfg ServerHTTP, handler http.Handler) *Server {
	srv := &http.Server{
		Addr:           net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
	}
	return &Server{
		srv: srv,
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
