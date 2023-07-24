package internalhttp

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type response struct {
	Message string `json:"message"`
}

func (h *Handler) newResponse(c *gin.Context, action string, code int, message string, err error) {
	if err != nil {
		h.Logger.Error(message, slog.String("action", action), slog.String("error", err.Error()))
	}
	c.AbortWithStatusJSON(code, response{Message: message})
}
