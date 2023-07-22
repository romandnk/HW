package internalhttp

import "net/http"

func (h *Handler) HelloWorld(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
