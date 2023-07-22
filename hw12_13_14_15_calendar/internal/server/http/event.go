package internalhttp

import (
	"github.com/gin-gonic/gin" //nolint:gci
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"net/http"
	"time"
)

type bodyEventCreate struct {
	Title                string `json:"title" binding:"required"`
	Date                 string `json:"date" binding:"required"`
	Duration             string `json:"duration" binding:"required"`
	Description          string `json:"description"`
	UserID               int    `json:"user_id" binding:"required"` //nolint:tagliatelle
	NotificationInterval string `json:"notification_interval"`      //nolint:tagliatelle
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var event models.Event
	var eventFromBody bodyEventCreate

	if err := c.ShouldBindJSON(&eventFromBody); err != nil {
		h.newResponse(c, "create", http.StatusBadRequest, "error parsing request body", err)
		return
	}

	date, err := time.Parse(time.RFC3339, eventFromBody.Date)
	if err != nil {
		h.newResponse(c, "create", http.StatusBadRequest, "error parsing date", err)
		return
	}
	duration, err := time.ParseDuration(eventFromBody.Duration)
	if err != nil {
		h.newResponse(c, "create", http.StatusBadRequest, "error parsing duration", err)
		return
	}

	var notificationInterval time.Duration
	if eventFromBody.NotificationInterval != "" {
		notificationInterval, err = time.ParseDuration(eventFromBody.NotificationInterval)
		if err != nil {
			h.newResponse(c, "create", http.StatusBadRequest, "error parsing notificationInterval", err)
			return
		}
	}

	event.Title = eventFromBody.Title
	event.Date = date
	event.Duration = duration
	event.Description = eventFromBody.Description
	event.UserID = eventFromBody.UserID
	event.NotificationInterval = notificationInterval

	id, err := h.Services.CreateEvent(c, event)
	if err != nil {
		h.newResponse(c, "create", http.StatusInternalServerError, "error creating event", err)
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) UpdateEvent(c *gin.Context) {
	return
}

func (h *Handler) DeleteEvent(c *gin.Context) {
	return
}

func (h *Handler) GetAllByDayEvents(c *gin.Context) {
	return
}

func (h *Handler) GetAllByWeekEvents(c *gin.Context) {
	return
}

func (h *Handler) GetAllByMonthEvents(c *gin.Context) {
	return
}
