package internalhttp

import (
	"context"
	"net/http"
	"time"
)

type ServerConf struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server struct {
	srv *http.Server
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application, cfg ServerConf, handler http.Handler) *Server {
	srv := &http.Server{
		Addr:           cfg.Host + ":" + cfg.Port,
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
