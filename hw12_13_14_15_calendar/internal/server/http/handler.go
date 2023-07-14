package internalhttp

import (
	sqlstorage "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/sql"
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

func (h *Handler) InitRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", h.HelloWorld)

	h.ServeMux = router

	return router
}
