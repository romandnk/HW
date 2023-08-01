package internalhttp

import (
	"errors"
	"github.com/gin-gonic/gin"
	customerror "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/errors"
	"golang.org/x/exp/slog"
)

type response struct {
	Action  string `json:"action"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Error   string `json:"errors"`
}

func newResponse(action, field, message string, err error) response {
	var customError customerror.CustomError

	if errors.As(err, &customError) {
		resp := response{
			Action:  action,
			Field:   customError.Field,
			Message: message,
			Error:   customError.Error(),
		}
		return resp
	}

	resp := response{
		Action:  action,
		Field:   field,
		Message: message,
		Error:   err.Error(),
	}

	return resp
}

func (h *HandlerHTTP) sentResponse(c *gin.Context, code int, resp response) {
	if resp.Error != "" {
		h.logger.Error(resp.Message,
			slog.String("action", resp.Action),
			slog.String("errors", resp.Error))
	}
	c.AbortWithStatusJSON(code, resp)
}
