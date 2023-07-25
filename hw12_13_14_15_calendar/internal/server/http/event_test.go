package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	mock_logger "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger/mock"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/service"
	"golang.org/x/exp/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	mock_service "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/service/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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

	handler := NewHandler(services, logger)

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
		name            string
		expectedErr     string
		expectedMessage string
		requestBody     map[string]interface{}
	}{
		{
			name:            "title is bool",
			expectedErr:     "json: cannot unmarshal bool into Go struct field bodyEvent.title of type string",
			expectedMessage: "error parsing request body",
			requestBody: map[string]interface{}{
				"title":    true,
				"date":     "2023-07-22T12:00:00Z",
				"duration": "1h30m",
				"user_id":  1,
			},
		},
		{
			name:            "invalid date",
			expectedErr:     "parsing time \"date\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"date\" as \"2006\"",
			expectedMessage: "error parsing date",
			requestBody: map[string]interface{}{
				"title":    "test",
				"date":     "date",
				"duration": "1h30m",
				"user_id":  1,
			},
		},
		{
			name:            "invalid duration",
			expectedErr:     "time: unknown unit \"y\" in duration \"1y1h30m\"",
			expectedMessage: "error parsing duration",
			requestBody: map[string]interface{}{
				"title":    "test",
				"date":     "2023-07-22T12:00:00Z",
				"duration": "1y1h30m",
				"user_id":  1,
			},
		},
		{
			name:            "invalid notificationInterval",
			expectedErr:     "time: invalid duration \"interval\"",
			expectedMessage: "error parsing notificationInterval",
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

			logger.EXPECT().Error(tc.expectedMessage, slog.String("action", "create"), slog.String("error", tc.expectedErr))

			handler := NewHandler(services, logger)

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
			require.Equal(t, tc.expectedMessage, message)
			require.True(t, ok)
		})
	}
}

//nolint:funlen
func TestHandlerCreateEventErrorCreatingEvent(t *testing.T) {
	testCases := []struct {
		name            string
		expectedErr     error
		expectedMessage string
		expectedEvent   models.Event
		requestBody     map[string]interface{}
	}{
		{
			name:            "title is empty",
			expectedErr:     service.ErrEmptyTitle,
			expectedMessage: "error creating event",
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
			name:            "user id is 0",
			expectedErr:     service.ErrInvalidUserID,
			expectedMessage: "error creating event",
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
			name:            "user id is -1",
			expectedErr:     service.ErrInvalidUserID,
			expectedMessage: "error creating event",
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
			name:            "duration is 0",
			expectedErr:     service.ErrInvalidDuration,
			expectedMessage: "error creating event",
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
			name:            "duration is -1 hour",
			expectedErr:     service.ErrInvalidDuration,
			expectedMessage: "error creating event",
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
			name:            "notification interval is -1 hour",
			expectedErr:     service.ErrInvalidDuration,
			expectedMessage: "error creating event",
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

			services.EXPECT().CreateEvent(gomock.Any(), tc.expectedEvent).Return("", tc.expectedErr)
			logger.EXPECT().Error(tc.expectedMessage,
				slog.String("action", "create"),
				slog.String("error", tc.expectedErr.Error()))

			handler := NewHandler(services, logger)

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
			require.Equal(t, tc.expectedMessage, message)
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

	handler := NewHandler(services, logger)

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
		name            string
		expectedErr     error
		expectedMessage string
		expectedEvent   models.Event
		requestBody     map[string]interface{}
	}{
		{
			name:            "duration is -1 hour",
			expectedErr:     service.ErrInvalidDuration,
			expectedMessage: "error updating event",
			expectedEvent: models.Event{
				Duration: -1 * time.Hour,
			},
			requestBody: map[string]interface{}{
				"duration": "-1h",
			},
		},
		{
			name:            "notification interval is -1 hour",
			expectedErr:     service.ErrInvalidDuration,
			expectedMessage: "error updating event",
			expectedEvent: models.Event{
				NotificationInterval: -time.Hour,
			},
			requestBody: map[string]interface{}{
				"notification_interval": "-1h",
			},
		},
		{
			name:            "user id is -1",
			expectedErr:     service.ErrInvalidUserID,
			expectedMessage: "error updating event",
			expectedEvent: models.Event{
				UserID: -1,
			},
			requestBody: map[string]interface{}{
				"user_id": -1,
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

			services.EXPECT().UpdateEvent(gomock.Any(), id, tc.expectedEvent).Return(models.Event{}, tc.expectedErr)
			logger.EXPECT().Error(tc.expectedMessage,
				slog.String("action", "update"),
				slog.String("error", tc.expectedErr.Error()))

			handler := NewHandler(services, logger)

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

			require.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			message, ok := responseBody["message"]
			require.Equal(t, tc.expectedMessage, message)
			require.True(t, ok)
		})
	}
}
