package internalhttp

import (
	"context"
	"net"
	"net/http"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/cmd/config"
)

type ServerHTTP struct {
	srv *http.Server
}

func NewServerHTTP(cfg config.ServerHTTPConfig, handler http.Handler) *ServerHTTP {
	srv := &http.Server{
		Addr:           net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
	}
	return &ServerHTTP{
		srv: srv,
	}
}

func (s *ServerHTTP) Start() error {
	return s.srv.ListenAndServe()
}

func (s *ServerHTTP) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
