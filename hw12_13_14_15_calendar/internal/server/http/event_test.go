package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mock_logger "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger/mock"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/service"
	mock_service "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/service/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/exp/slog"
)

const url = "/api/v1/events"

func TestHandlerCreateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)
	logger := mock_logger.NewMockLogger(ctrl)

	expectedEvent := models.Event{
		Title:                "Test Event",
		Date:                 time.Date(2023, 7, 22, 12, 0, 0, 0, time.UTC),
		Duration:             1*time.Hour + 30*time.Minute,
		Description:          "This is a test event",
		UserID:               1,
		NotificationInterval: 10 * time.Minute,
	}
	expectedID := "test uuid"
	services.EXPECT().CreateEvent(gomock.Any(), expectedEvent).Return(expectedID, nil)

	handler := NewHandlerHTTP(services, logger)

	r := gin.Default()
	r.POST(url, handler.CreateEvent)

	requestBody := map[string]interface{}{
		"title":                 "Test Event",
		"date":                  "2023-07-22T12:00:00Z",
		"duration":              "1h30m",
		"description":           "This is a test event",
		"user_id":               1,
		"notification_interval": "10m",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	id, ok := responseBody["id"]
	require.Equal(t, expectedID, id)
	require.True(t, ok)
}

func TestHandlerCreateEventError(t *testing.T) {
	testCases := []struct {
		name             string
		expectedResponse response
		requestBody      map[string]interface{}
	}{
		{
			name: "title is bool",
			expectedResponse: response{
				Action:  createAction,
				Message: ErrParsingBody.Error(),
				Error:   "json: cannot unmarshal bool into Go struct field bodyEvent.title of type string",
			},
			requestBody: map[string]interface{}{
				"title":    true,
				"date":     "2023-07-22T12:00:00Z",
				"duration": "1h30m",
				"user_id":  1,
			},
		},
		{
			name: "invalid date",
			expectedResponse: response{
				Action:  createAction,
				Message: ErrParsingDate.Error(),
				Error:   "parsing time \"date\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"date\" as \"2006\"",
			},
			requestBody: map[string]interface{}{
				"title":    "test",
				"date":     "date",
				"duration": "1h30m",
				"user_id":  1,
			},
		},
		{
			name: "invalid duration",
			expectedResponse: response{
				Action:  createAction,
				Message: ErrParsingDuration.Error(),
				Error:   "time: unknown unit \"y\" in duration \"1y1h30m\"",
			},
			requestBody: map[string]interface{}{
				"title":    "test",
				"date":     "2023-07-22T12:00:00Z",
				"duration": "1y1h30m",
				"user_id":  1,
			},
		},
		{
			name: "invalid notification_interval",
			expectedResponse: response{
				Action:  createAction,
				Message: ErrParsingNotificationInterval.Error(),
				Error:   "time: invalid duration \"interval\"",
			},
			requestBody: map[string]interface{}{
				"title":                 "test",
				"date":                  "2023-07-22T12:00:00Z",
				"duration":              "1h30m",
				"user_id":               1,
				"notification_interval": "interval",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			services := mock_service.NewMockServices(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			logger.EXPECT().Error(tc.expectedResponse.Message,
				slog.String("action", "create"),
				slog.String("errors", tc.expectedResponse.Error))

			handler := NewHandlerHTTP(services, logger)

			r := gin.Default()
			r.POST(url, handler.CreateEvent)

			jsonBody, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			message, ok := responseBody["message"]
			require.Equal(t, tc.expectedResponse.Message, message)
			require.True(t, ok)
		})
	}
}

//nolint:funlen
func TestHandlerCreateEventErrorCreatingEvent(t *testing.T) {
	testCases := []struct {
		name             string
		expectedResponse response
		expectedEvent    models.Event
		requestBody      map[string]interface{}
	}{
		{
			name: "title is empty",
			expectedResponse: response{
				Action:  createAction,
				Field:   "title",
				Message: "error creating event",
				Error:   "title cannot be empty",
			},
			expectedEvent: models.Event{
				Title:    "",
				Date:     time.Date(2023, 7, 22, 12, 0, 0, 0, time.UTC),
				Duration: 1*time.Hour + 30*time.Minute,
				UserID:   1,
			},
			requestBody: map[string]interface{}{
				"title":    "",
				"date":     "2023-07-22T12:00:00Z",
				"duration": "1h30m",
				"user_id":  1,
			},
		},
		{
			name: "user id is 0",
			expectedResponse: response{
				Action:  createAction,
				Field:   "user_id",
				Message: "error creating event",
				Error:   "user id must not be positive number",
			},
			expectedEvent: models.Event{
				Title:    "test",
				Date:     time.Date(2023, 7, 22, 12, 0, 0, 0, time.UTC),
				Duration: 1*time.Hour + 30*time.Minute,
				UserID:   0,
			},
			requestBody: map[string]interface{}{
				"title":    "test",
				"date":     "2023-07-22T12:00:00Z",
				"duration": "1h30m",
				"user_id":  0,
			},
		},
		{
			name: "user id is -1",
			expectedResponse: response{
				Action:  createAction,
				Field:   "user_id",
				Message: "error creating event",
				Error:   "user id must not be positive number",
			},
			expectedEvent: models.Event{
				Title:    "test",
				Date:     time.Date(2023, 7, 22, 12, 0, 0, 0, time.UTC),
				Duration: 1*time.Hour + 30*time.Minute,
				UserID:   -1,
			},
			requestBody: map[string]interface{}{
				"title":    "test",
				"date":     "2023-07-22T12:00:00Z",
				"duration": "1h30m",
				"user_id":  -1,
			},
		},
		{
			name: "duration is 0",
			expectedResponse: response{
				Action:  createAction,
				Field:   "duration",
				Message: "error creating event",
				Error:   "duration cannot be non-positive",
			},
			expectedEvent: models.Event{
				Title:    "test",
				Date:     time.Date(2023, 7, 22, 12, 0, 0, 0, time.UTC),
				Duration: 0,
				UserID:   1,
			},
			requestBody: map[string]interface{}{
				"title":    "test",
				"date":     "2023-07-22T12:00:00Z",
				"duration": "0s",
				"user_id":  1,
			},
		},
		{
			name: "duration is -1 hour",
			expectedResponse: response{
				Action:  createAction,
				Field:   "duration",
				Message: "error creating event",
				Error:   "duration cannot be non-positive",
			},
			expectedEvent: models.Event{
				Title:    "test",
				Date:     time.Date(2023, 7, 22, 12, 0, 0, 0, time.UTC),
				Duration: -1 * time.Hour,
				UserID:   1,
			},
			requestBody: map[string]interface{}{
				"title":    "test",
				"date":     "2023-07-22T12:00:00Z",
				"duration": "-1h",
				"user_id":  1,
			},
		},
		{
			name: "notification interval is -1 hour",
			expectedResponse: response{
				Action:  createAction,
				Field:   "notification_interval",
				Message: "error creating event",
				Error:   "notification interval cannot be negative",
			},
			expectedEvent: models.Event{
				Title:                "test",
				Date:                 time.Date(2023, 7, 22, 12, 0, 0, 0, time.UTC),
				Duration:             1 * time.Hour,
				UserID:               1,
				NotificationInterval: -time.Hour,
			},
			requestBody: map[string]interface{}{
				"title":                 "test",
				"date":                  "2023-07-22T12:00:00Z",
				"duration":              "1h",
				"user_id":               1,
				"notification_interval": "-1h",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			services := mock_service.NewMockServices(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			services.EXPECT().CreateEvent(gomock.Any(), tc.expectedEvent).Return("", errors.New(tc.expectedResponse.Error))
			logger.EXPECT().Error(tc.expectedResponse.Message,
				slog.String("action", tc.expectedResponse.Action),
				slog.String("errors", tc.expectedResponse.Error))

			handler := NewHandlerHTTP(services, logger)

			r := gin.Default()
			r.POST(url, handler.CreateEvent)

			jsonBody, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusInternalServerError, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			message, ok := responseBody["message"]
			require.Equal(t, tc.expectedResponse.Message, message)
			require.True(t, ok)
		})
	}
}

func TestHandlerUpdateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)
	logger := mock_logger.NewMockLogger(ctrl)

	id := uuid.New().String()

	event := models.Event{
		Title:                "Test Event update",
		Date:                 time.Date(2023, 7, 22, 12, 0, 0, 0, time.UTC),
		Duration:             1*time.Hour + 30*time.Minute,
		Description:          "This is a test event update",
		UserID:               1,
		NotificationInterval: 10 * time.Minute,
	}

	expectedEvent := event
	expectedEvent.ID = id

	services.EXPECT().UpdateEvent(gomock.Any(), id, event).Return(expectedEvent, nil)

	handler := NewHandlerHTTP(services, logger)

	r := gin.Default()
	r.PATCH(url+"/:id", handler.UpdateEvent)

	requestBody := map[string]interface{}{
		"id":                    id,
		"title":                 "Test Event update",
		"date":                  "2023-07-22T12:00:00Z",
		"duration":              "1h30m",
		"description":           "This is a test event update",
		"user_id":               1,
		"notification_interval": "10m",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url+"/"+id, bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var responseBody Response
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	expectedBody := Response{
		ID:                   id,
		Title:                "Test Event update",
		Date:                 "2023-07-22T12:00:00Z",
		Duration:             "1h30m0s",
		Description:          "This is a test event update",
		UserID:               1,
		NotificationInterval: "10m0s",
	}

	require.Equal(t, expectedBody, responseBody)
}

func TestHandlerUpdateEventError(t *testing.T) {
	testCases := []struct {
		name             string
		expectedResponse response
		expectedEvent    models.Event
		id               string
		requestBody      map[string]interface{}
	}{
		{
			name: "incorrect id",
			expectedResponse: response{
				Action:  updateAction,
				Field:   "id (param)",
				Message: ErrInvalidID.Error(),
				Error:   "invalid UUID length: 7",
			},
			id: "1234567",
		},
		{
			name: "invalid request body",
			expectedResponse: response{
				Action:  updateAction,
				Field:   "",
				Message: ErrParsingBody.Error(),
				Error:   "json: cannot unmarshal bool into Go struct field bodyEvent.title of type string",
			},
			id: uuid.New().String(),
			requestBody: map[string]interface{}{
				"title":                 true,
				"date":                  "2023-07-22T12:00:00Z",
				"duration":              "1h30m",
				"description":           "This is a test event update",
				"user_id":               1,
				"notification_interval": "10m",
			},
		},
		{
			name: "invalid date",
			expectedResponse: response{
				Action:  updateAction,
				Field:   "date",
				Message: ErrParsingDate.Error(),
				Error:   "parsing time \"date\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"date\" as \"2006\"",
			},
			id: uuid.New().String(),
			requestBody: map[string]interface{}{
				"title":                 "test title",
				"date":                  "date",
				"duration":              "1h30m",
				"description":           "This is a test event update",
				"user_id":               1,
				"notification_interval": "10m",
			},
		},
		{
			name: "invalid duration",
			expectedResponse: response{
				Action:  updateAction,
				Field:   "duration",
				Message: ErrParsingDuration.Error(),
				Error:   "time: invalid duration \"interval\"",
			},
			id: uuid.New().String(),
			requestBody: map[string]interface{}{
				"title":                 "test title",
				"date":                  "2023-07-22T12:00:00Z",
				"duration":              "duration",
				"description":           "This is a test event update",
				"user_id":               1,
				"notification_interval": "10m",
			},
		},
		{
			name: "invalid notification interval",
			expectedResponse: response{
				Action:  updateAction,
				Field:   "notification_interval",
				Message: ErrParsingNotificationInterval.Error(),
				Error:   "time: invalid duration \"interval\"",
			},
			id: uuid.New().String(),
			requestBody: map[string]interface{}{
				"title":                 "test title",
				"date":                  "2023-07-22T12:00:00Z",
				"duration":              "1h30m",
				"description":           "This is a test event update",
				"user_id":               1,
				"notification_interval": "interval",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			services := mock_service.NewMockServices(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			logger.EXPECT().Error(tc.expectedResponse.Message,
				slog.String("action", "update"),
				slog.String("errors", tc.expectedResponse.Error))

			handler := NewHandlerHTTP(services, logger)

			r := gin.Default()
			r.PATCH(url+"/:id", handler.UpdateEvent)

			jsonBody, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url+"/"+tc.id, bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			message, ok := responseBody["message"]
			require.Equal(t, tc.expectedResponse.Message, message)
			require.True(t, ok)
		})
	}
}

func TestHandlerUpdateEventErrorUpdatingEvent(t *testing.T) {
	testCases := []struct {
		name             string
		expectedResponse response
		expectedEvent    models.Event
		requestBody      map[string]interface{}
	}{
		{
			name: "duration is -1 hour",
			expectedResponse: response{
				Action:  updateAction,
				Field:   "duration",
				Message: "error updating event",
				Error:   service.ErrInvalidDuration.Error(),
			},
			expectedEvent: models.Event{
				Duration: -1 * time.Hour,
			},
			requestBody: map[string]interface{}{
				"duration": "-1h",
			},
		},
		{
			name: "notification interval is -1 hour",
			expectedResponse: response{
				Action:  updateAction,
				Field:   "notification_interval",
				Message: "error updating event",
				Error:   service.ErrInvalidDuration.Error(),
			},
			expectedEvent: models.Event{
				NotificationInterval: -time.Hour,
			},
			requestBody: map[string]interface{}{
				"notification_interval": "-1h",
			},
		},
		{
			name: "user id is -1",
			expectedResponse: response{
				Action:  updateAction,
				Field:   "notification_interval",
				Message: "error updating event",
				Error:   service.ErrInvalidUserID.Error(),
			},
			expectedEvent: models.Event{
				UserID: -1,
			},
			requestBody: map[string]interface{}{
				"user_id": -1,
			},
		},
		{
			name: "title contains only spaces",
			expectedResponse: response{
				Action:  updateAction,
				Field:   "notification_interval",
				Message: "error updating event",
				Error:   service.ErrEmptyTitle.Error(),
			},
			expectedEvent: models.Event{
				Title: "                 ",
			},
			requestBody: map[string]interface{}{
				"title": "                 ",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			services := mock_service.NewMockServices(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			id := uuid.New().String()

			services.EXPECT().UpdateEvent(gomock.Any(), id, tc.expectedEvent).
				Return(models.Event{}, errors.New(tc.expectedResponse.Error))
			logger.EXPECT().Error(tc.expectedResponse.Message,
				slog.String("action", "update"),
				slog.String("errors", tc.expectedResponse.Error))

			handler := NewHandlerHTTP(services, logger)

			r := gin.Default()
			r.PATCH(url+"/:id", handler.UpdateEvent)

			jsonBody, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url+"/"+id, bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusInternalServerError, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			message, ok := responseBody["message"]
			require.Equal(t, tc.expectedResponse.Message, message)
			require.True(t, ok)
		})
	}
}

func TestHandlerDeleteEvent(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)
	logger := mock_logger.NewMockLogger(ctrl)

	id := uuid.New().String()

	services.EXPECT().DeleteEvent(gomock.Any(), id).Return(nil)

	handler := NewHandlerHTTP(services, logger)

	r := gin.Default()
	r.DELETE(url+"/:id", handler.DeleteEvent)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url+"/"+id, nil)
	require.NoError(t, err)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandlerDeleteEventError(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)
	logger := mock_logger.NewMockLogger(ctrl)

	id := "test id"

	expectedMessage := "invalid id"
	expectedError := "invalid UUID length: 7"

	logger.EXPECT().Error(expectedMessage,
		slog.String("action", "delete"),
		slog.String("errors", expectedError))

	handler := NewHandlerHTTP(services, logger)

	r := gin.Default()
	r.DELETE(url+"/:id", handler.DeleteEvent)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url+"/"+id, nil)
	require.NoError(t, err)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerDeleteEventErrorDeletingEvent(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)
	logger := mock_logger.NewMockLogger(ctrl)

	id := uuid.New().String()

	expectedResponse := response{
		Action:  deleteAction,
		Field:   "id (param)",
		Message: "error deleting event",
		Error:   "no event with id" + id,
	}

	services.EXPECT().DeleteEvent(gomock.Any(), id).Return(errors.New(expectedResponse.Error))
	logger.EXPECT().Error(expectedResponse.Message,
		slog.String("action", "delete"),
		slog.String("errors", expectedResponse.Error))

	handler := NewHandlerHTTP(services, logger)

	r := gin.Default()
	r.DELETE(url+"/:id", handler.DeleteEvent)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url+"/"+id, nil)
	require.NoError(t, err)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandlerGetAllEvents(t *testing.T) {
	dateStr := "2023-07-22T12:00:00Z"
	date, err := time.Parse(time.RFC3339, dateStr)
	require.NoError(t, err)

	testCases := []struct {
		name           string
		period         string
		expectedEvents []models.Event
	}{
		{
			name:   "get all by day",
			period: "day",
			expectedEvents: []models.Event{
				{
					ID:                   uuid.New().String(),
					Title:                "1",
					Date:                 date,
					Duration:             time.Second,
					Description:          "1",
					UserID:               1,
					NotificationInterval: time.Second,
				},
				{
					ID:                   uuid.New().String(),
					Title:                "2",
					Date:                 date,
					Duration:             2 * time.Second,
					Description:          "2",
					UserID:               2,
					NotificationInterval: 2 * time.Second,
				},
				{
					ID:                   uuid.New().String(),
					Title:                "6",
					Date:                 date,
					Duration:             6 * time.Second,
					Description:          "6",
					UserID:               6,
					NotificationInterval: 6 * time.Second,
				},
			},
		},
		{
			name:   "get all by week",
			period: "week",
			expectedEvents: []models.Event{
				{
					ID:                   uuid.New().String(),
					Title:                "1",
					Date:                 date,
					Duration:             time.Second,
					Description:          "1",
					UserID:               1,
					NotificationInterval: time.Second,
				},
				{
					ID:                   uuid.New().String(),
					Title:                "2",
					Date:                 date.AddDate(0, 0, 2),
					Duration:             2 * time.Second,
					Description:          "2",
					UserID:               2,
					NotificationInterval: 2 * time.Second,
				},
				{
					ID:                   uuid.New().String(),
					Title:                "6",
					Date:                 date.AddDate(0, 0, 6),
					Duration:             6 * time.Second,
					Description:          "6",
					UserID:               6,
					NotificationInterval: 6 * time.Second,
				},
			},
		},
		{
			name:   "get all by month",
			period: "month",
			expectedEvents: []models.Event{
				{
					ID:                   uuid.New().String(),
					Title:                "1",
					Date:                 date,
					Duration:             time.Second,
					Description:          "1",
					UserID:               1,
					NotificationInterval: time.Second,
				},
				{
					ID:                   uuid.New().String(),
					Title:                "2",
					Date:                 date.AddDate(0, 0, 15),
					Duration:             2 * time.Second,
					Description:          "2",
					UserID:               2,
					NotificationInterval: 2 * time.Second,
				},
				{
					ID:                   uuid.New().String(),
					Title:                "6",
					Date:                 date.AddDate(0, 0, 29),
					Duration:             6 * time.Second,
					Description:          "6",
					UserID:               6,
					NotificationInterval: 6 * time.Second,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			services := mock_service.NewMockServices(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			handler := NewHandlerHTTP(services, logger)

			r := gin.Default()

			w := httptest.NewRecorder()

			switch tc.period {
			case "day":
				services.EXPECT().GetAllByDayEvents(gomock.Any(), date).Return(tc.expectedEvents, nil)
				r.GET(url+"/"+tc.period+"/:date", handler.GetAllByDayEvents)
			case "week":
				services.EXPECT().GetAllByWeekEvents(gomock.Any(), date).Return(tc.expectedEvents, nil)
				r.GET(url+"/"+tc.period+"/:date", handler.GetAllByWeekEvents)
			case "month":
				services.EXPECT().GetAllByMonthEvents(gomock.Any(), date).Return(tc.expectedEvents, nil)
				r.GET(url+"/"+tc.period+"/:date", handler.GetAllByMonthEvents)
			}

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/"+tc.period+"/"+dateStr, nil)
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)

			var responseBody eventsResponse
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			require.Equal(t, 3, responseBody.Total)
			require.Equal(t, formResponseGetBy(tc.expectedEvents).Data, responseBody.Data)
		})
	}
}

func TestHandlerGetAllEventsError(t *testing.T) {
	date := "date"
	expectedMessage := "date must be in RFC3339 format"
	expectedError := "parsing time \"date\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"date\" as \"2006\""

	testCases := []struct {
		name   string
		period string
		action string
	}{
		{
			name:   "error parsing date by day",
			period: "day",
			action: "get by day",
		},
		{
			name:   "error parsing date by week",
			period: "week",
			action: "get by week",
		},
		{
			name:   "error parsing date by month",
			period: "month",
			action: "get by month",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			services := mock_service.NewMockServices(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			logger.EXPECT().Error(expectedMessage,
				slog.String("action", tc.action),
				slog.String("errors", expectedError))

			handler := NewHandlerHTTP(services, logger)

			r := gin.Default()

			w := httptest.NewRecorder()

			switch tc.period {
			case "day":
				r.GET(url+"/"+tc.period+"/:date", handler.GetAllByDayEvents)
			case "week":
				r.GET(url+"/"+tc.period+"/:date", handler.GetAllByWeekEvents)
			case "month":
				r.GET(url+"/"+tc.period+"/:date", handler.GetAllByMonthEvents)
			}

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/"+tc.period+"/"+date, nil)
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}
