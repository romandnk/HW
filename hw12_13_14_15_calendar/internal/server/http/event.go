package internalhttp

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

type bodyEvent struct {
	Title                string `json:"title"`
	Date                 string `json:"date"`
	Duration             string `json:"duration"`
	Description          string `json:"description"`
	UserID               int    `json:"user_id"`
	NotificationInterval string `json:"notification_interval"`
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var event models.Event
	var eventFromBody bodyEvent

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

	id, err := h.services.CreateEvent(c, event)
	if err != nil {
		h.newResponse(c, "create", http.StatusBadRequest, "error creating event", err)
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

type Response struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Date                 string `json:"date"`
	Duration             string `json:"duration"`
	Description          string `json:"description"`
	UserID               int    `json:"user_id"`
	NotificationInterval string `json:"notification_interval"`
}

func (h *Handler) UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.newResponse(c, "update", http.StatusBadRequest, "invalid id", err)
		return
	}

	var eventFromBody bodyEvent

	if err := c.ShouldBindJSON(&eventFromBody); err != nil {
		h.newResponse(c, "update", http.StatusBadRequest, "error parsing request body", err)
		return
	}

	var date time.Time
	if eventFromBody.Date != "" {
		date, err = time.Parse(time.RFC3339, eventFromBody.Date)
		if err != nil {
			h.newResponse(c, "update", http.StatusBadRequest, "error parsing date", err)
			return
		}
	}

	var duration time.Duration
	if eventFromBody.Duration != "" {
		duration, err = time.ParseDuration(eventFromBody.Duration)
		if err != nil {
			h.newResponse(c, "update", http.StatusBadRequest, "error parsing duration", err)
			return
		}
	}

	var notificationInterval time.Duration
	if eventFromBody.NotificationInterval != "" {
		notificationInterval, err = time.ParseDuration(eventFromBody.NotificationInterval)
		if err != nil {
			h.newResponse(c, "update", http.StatusBadRequest, "error parsing notification interval", err)
			return
		}
	}

	event := models.Event{
		Title:                eventFromBody.Title,
		Date:                 date,
		Duration:             duration,
		Description:          eventFromBody.Description,
		UserID:               eventFromBody.UserID,
		NotificationInterval: notificationInterval,
	}

	updatedEvent, err := h.services.UpdateEvent(c, parsedID.String(), event)
	if err != nil {
		h.newResponse(c, "update", http.StatusBadRequest, "error updating event", err)
		return
	}

	c.JSON(http.StatusOK, Response{
		ID:                   updatedEvent.ID,
		Title:                updatedEvent.Title,
		Date:                 updatedEvent.Date.Format(time.RFC3339),
		Duration:             updatedEvent.Duration.String(),
		Description:          updatedEvent.Description,
		UserID:               updatedEvent.UserID,
		NotificationInterval: updatedEvent.NotificationInterval.String(),
	})
}

func (h *Handler) DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.newResponse(c, "delete", http.StatusBadRequest, "invalid id", err)
		return
	}

	err = h.services.DeleteEvent(c, parsedID.String())
	if err != nil {
		h.newResponse(c, "delete", http.StatusBadRequest, "error deleting event", err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetAllByDayEvents(c *gin.Context) {
	date := c.Param("date")
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		h.newResponse(c, "get by day", http.StatusBadRequest, "error parsing date", err)
		return
	}

	events, err := h.services.GetAllByDayEvents(c, parsedDate)
	if err != nil {
		h.newResponse(c, "get by day", http.StatusBadRequest, "error getting events by day", err)
		return
	}

	c.JSON(http.StatusOK, struct {
		Total int            `json:"total"`
		Data  []models.Event `json:"data"`
	}{
		Total: len(events),
		Data:  events,
	})
}

func (h *Handler) GetAllByWeekEvents(c *gin.Context) {
	date := c.Param("date")
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		h.newResponse(c, "get by week", http.StatusBadRequest, "error parsing date", err)
		return
	}

	events, err := h.services.GetAllByWeekEvents(c, parsedDate)
	if err != nil {
		h.newResponse(c, "get by week", http.StatusBadRequest, "error getting events by week", err)
		return
	}

	c.JSON(http.StatusOK, struct {
		Total int            `json:"total"`
		Data  []models.Event `json:"data"`
	}{
		Total: len(events),
		Data:  events,
	})
}

func (h *Handler) GetAllByMonthEvents(c *gin.Context) {
	date := c.Param("date")
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		h.newResponse(c, "get by month", http.StatusBadRequest, "error parsing date", err)
		return
	}

	events, err := h.services.GetAllByMonthEvents(c, parsedDate)
	if err != nil {
		h.newResponse(c, "get by month", http.StatusBadRequest, "error getting events by month", err)
		return
	}

	c.JSON(http.StatusOK, struct {
		Total int            `json:"total"`
		Data  []models.Event `json:"data"`
	}{
		Total: len(events),
		Data:  events,
	})
}
