package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestHandlerCreateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)
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

	handler := NewHandler(services)

	url := "/api/v1/events"
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
