package internalhttp

import (
	"net/http"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/service"
)

type Handler struct {
	*http.ServeMux
	Services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		Services: services,
	}
}

func (h *Handler) InitRoutes(log *logger.MyLogger) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", middlewareLogging(log, h.HelloWorld))

	h.ServeMux = router

	return router
}
