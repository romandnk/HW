package internalhttp

import (
	sqlstorage "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/sql"
	"golang.org/x/exp/slog"
	"net/http"
)

type Handler struct {
	*http.ServeMux
	Event sqlstorage.Event
}

func NewHandler(event sqlstorage.Event) *Handler {
	return &Handler{
		Event: event,
	}
}

func (h *Handler) InitRoutes(log *slog.Logger) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", middlewareLogging(log, h.HelloWorld))

	h.ServeMux = router

	return router
}
