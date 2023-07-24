package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	mock_logger "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger/mock"
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
			name:            "empty title",
			expectedErr:     "Key: 'bodyEventCreate.Title' Error:Field validation for 'Title' failed on the 'required' tag",
			expectedMessage: "error parsing request body",
			requestBody: map[string]interface{}{
				"title":    "",
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

func TestHandlerCreateEventErrorCreatingEvent(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)
	logger := mock_logger.NewMockLogger(ctrl)

	expectedMessage := "error creating event"
	expectedErr := errors.New("user id must be positive number")

	expectedEvent := models.Event{
		Title:                "Test Event",
		Date:                 time.Date(2023, 7, 22, 12, 0, 0, 0, time.UTC),
		Duration:             1*time.Hour + 30*time.Minute,
		Description:          "This is a test event",
		UserID:               -1,
		NotificationInterval: 10 * time.Minute,
	}

	services.EXPECT().CreateEvent(gomock.Any(), expectedEvent).Return("", expectedErr)
	logger.EXPECT().Error(expectedMessage, slog.String("action", "create"), slog.String("error", expectedErr.Error()))

	handler := NewHandler(services, logger)

	r := gin.Default()
	r.POST(url, handler.CreateEvent)

	requestBody := map[string]interface{}{
		"title":                 "Test Event",
		"date":                  "2023-07-22T12:00:00Z",
		"duration":              "1h30m",
		"description":           "This is a test event",
		"user_id":               -1,
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

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	message, ok := responseBody["message"]
	require.Equal(t, expectedMessage, message)
	require.True(t, ok)
}

//func TestHandlerUpdateEvent(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	services := mock_service.NewMockServices(ctrl)
//	logger := mock_logger.NewMockLogger(ctrl)
//
//	expectedEvent := models.Event{
//		ID:                   "test uuid",
//		Title:                "Test Event update",
//		Date:                 time.Date(2023, 7, 22, 12, 0, 0, 0, time.UTC),
//		Duration:             1*time.Hour + 30*time.Minute,
//		Description:          "This is a test event update",
//		UserID:               1,
//		NotificationInterval: 10 * time.Minute,
//	}
//
//	services.EXPECT().UpdateEvent(gomock.Any(), expectedEvent.ID, expectedEvent).Return(expectedEvent, nil)
//
//	handler := NewHandler(services, logger)
//
//	r := gin.Default()
//	r.PATCH(url, handler.UpdateEvent)
//
//	requestBody := map[string]interface{}{
//		"title":                 "Test Event update",
//		"date":                  "2023-07-22T12:00:00Z",
//		"duration":              "1h30m",
//		"description":           "This is a test event update",
//		"user_id":               1,
//		"notification_interval": "10m",
//	}
//
//	jsonBody, err := json.Marshal(requestBody)
//	require.NoError(t, err)
//
//	w := httptest.NewRecorder()
//
//	ctx := context.Background()
//	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(jsonBody))
//	require.NoError(t, err)
//	req.Header.Set("Content-Type", "application/json")
//
//	r.ServeHTTP(w, req)
//
//	require.Equal(t, http.StatusOK, w.Code)
//
//	var responseBody map[string]interface{}
//	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
//	require.NoError(t, err)
//
//	expectedBody := map[string]interface{}{
//		"id":                    "test uuid",
//		"title":                 "Test Event update",
//		"date":                  "2023-07-22T12:00:00Z",
//		"duration":              "1h30m",
//		"description":           "This is a test event update",
//		"user_id":               1,
//		"notification_interval": "10m",
//	}
//	jsonExpectedBody, err := json.Marshal(&expectedBody)
//	require.NoError(t, err)
//	require.Equal(t, jsonExpectedBody, responseBody)
//}
