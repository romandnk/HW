package internalhttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
)

var (
	createAction     = "create"
	updateAction     = "update"
	deleteAction     = "delete"
	getByDayAction   = "get by day"
	getByWeekAction  = "get by week"
	getByMonthAction = "get by month"
)

var (
	ErrParsingDate                 = errors.New("date must be in RFC3339 format")
	ErrParsingBody                 = errors.New("error parsing json body")
	ErrParsingDuration             = errors.New("duration must be represented only in hours, minutes, seconds")
	ErrParsingNotificationInterval = errors.New("notification_interval must be represented only in hours, minutes, seconds")
	ErrInvalidID                   = errors.New("invalid id")
)

type bodyEvent struct {
	Title                string `json:"title"`
	Date                 string `json:"date"`
	Duration             string `json:"duration"`
	Description          string `json:"description"`
	UserID               int    `json:"user_id"`
	NotificationInterval string `json:"notification_interval"`
}

func (h *HandlerHTTP) CreateEvent(c *gin.Context) {
	var event models.Event
	var eventFromBody bodyEvent

	if err := c.ShouldBindJSON(&eventFromBody); err != nil {
		resp := newResponse(createAction, "", ErrParsingBody.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	date, err := time.Parse(time.RFC3339, eventFromBody.Date)
	if err != nil {
		resp := newResponse(createAction, "date", ErrParsingDate.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}
	duration, err := time.ParseDuration(eventFromBody.Duration)
	if err != nil {
		resp := newResponse(createAction, "duration", ErrParsingDuration.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	var notificationInterval time.Duration
	if eventFromBody.NotificationInterval != "" {
		notificationInterval, err = time.ParseDuration(eventFromBody.NotificationInterval)
		if err != nil {
			resp := newResponse(createAction, "notification_interval", ErrParsingNotificationInterval.Error(), err)
			h.sentResponse(c, http.StatusBadRequest, resp)
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
		message := "error creating event"
		resp := newResponse(createAction, "", message, err)
		h.sentResponse(c, http.StatusInternalServerError, resp)
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

func (h *HandlerHTTP) UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		resp := newResponse(updateAction, "id (param)", ErrInvalidID.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	var eventFromBody bodyEvent

	if err := c.ShouldBindJSON(&eventFromBody); err != nil {
		resp := newResponse(updateAction, "", ErrParsingBody.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	var date time.Time
	if eventFromBody.Date != "" {
		date, err = time.Parse(time.RFC3339, eventFromBody.Date)
		if err != nil {
			resp := newResponse(updateAction, "date", ErrParsingDate.Error(), err)
			h.sentResponse(c, http.StatusBadRequest, resp)
			return
		}
	}

	var duration time.Duration
	if eventFromBody.Duration != "" {
		duration, err = time.ParseDuration(eventFromBody.Duration)
		if err != nil {
			resp := newResponse(updateAction, "duration", ErrParsingDuration.Error(), err)
			h.sentResponse(c, http.StatusBadRequest, resp)
			return
		}
	}

	var notificationInterval time.Duration
	if eventFromBody.NotificationInterval != "" {
		notificationInterval, err = time.ParseDuration(eventFromBody.NotificationInterval)
		if err != nil {
			resp := newResponse(updateAction, "notification_interval", ErrParsingNotificationInterval.Error(), err)
			h.sentResponse(c, http.StatusBadRequest, resp)
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
		message := "error updating event"
		resp := newResponse(updateAction, "", message, err)
		h.sentResponse(c, http.StatusInternalServerError, resp)
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

func (h *HandlerHTTP) DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		resp := newResponse(deleteAction, "id (param)", ErrInvalidID.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	err = h.services.DeleteEvent(c, parsedID.String())
	if err != nil {
		message := "error deleting event"
		resp := newResponse(deleteAction, "", message, err)
		h.sentResponse(c, http.StatusInternalServerError, resp)
		return
	}

	c.Status(http.StatusOK)
}

type eventsResponse struct {
	Total int            `json:"total"`
	Data  []eventDetails `json:"data"`
}

type eventDetails struct {
	ID                   string        `json:"id"`
	Title                string        `json:"title"`
	Date                 time.Time     `json:"date"`
	Duration             time.Duration `json:"duration"`
	Description          string        `json:"description"`
	UserID               int           `json:"user_id"`
	NotificationInterval time.Duration `json:"notification_interval"`
}

func (h *HandlerHTTP) GetAllByDayEvents(c *gin.Context) {
	date := c.Param("date")
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		resp := newResponse(getByDayAction, "date (param)", ErrParsingDate.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	events, err := h.services.GetAllByDayEvents(c, parsedDate)
	if err != nil {
		message := "error getting events by day"
		resp := newResponse(getByDayAction, "", message, err)
		h.sentResponse(c, http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, formResponseGetBy(events))
}

func (h *HandlerHTTP) GetAllByWeekEvents(c *gin.Context) {
	date := c.Param("date")
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		resp := newResponse(getByWeekAction, "date (param)", ErrParsingDate.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	events, err := h.services.GetAllByWeekEvents(c, parsedDate)
	if err != nil {
		message := "error getting events by week"
		resp := newResponse(getByWeekAction, "", message, err)
		h.sentResponse(c, http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, formResponseGetBy(events))
}

func (h *HandlerHTTP) GetAllByMonthEvents(c *gin.Context) {
	date := c.Param("date")
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		resp := newResponse(getByMonthAction, "date (param)", ErrParsingDate.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	events, err := h.services.GetAllByMonthEvents(c, parsedDate)
	if err != nil {
		message := "error getting events by month"
		resp := newResponse(getByMonthAction, "", message, err)
		h.sentResponse(c, http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, formResponseGetBy(events))
}

func formResponseGetBy(events []models.Event) eventsResponse {
	var response eventsResponse
	response.Total = len(events)
	for _, event := range events {
		response.Data = append(response.Data, eventDetails{
			ID:                   event.ID,
			Title:                event.Title,
			Date:                 event.Date,
			Duration:             event.Duration,
			Description:          event.Description,
			UserID:               event.UserID,
			NotificationInterval: event.NotificationInterval,
		})
	}
	return response
}
