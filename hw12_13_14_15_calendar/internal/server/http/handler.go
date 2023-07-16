package internalhttp

import (
	"net/http"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage"
	"golang.org/x/exp/slog"
)

type Handler struct {
	*http.ServeMux
	Storage storage.StoreEvent
}

func NewHandler(storage storage.StoreEvent) *Handler {
	return &Handler{
		Storage: storage,
	}
}

func (h *Handler) InitRoutes(log *slog.Logger) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", middlewareLogging(log, h.HelloWorld))

	h.ServeMux = router

	return router
}
