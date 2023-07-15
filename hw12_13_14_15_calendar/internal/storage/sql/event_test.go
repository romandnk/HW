package sqlstorage

import (
	"context"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	mockdb "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/sql/.mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestStorageCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStoreEvent(ctrl)

	event := models.Event{
		Title:                "test title",
		Date:                 time.Now(),
		Duration:             time.Second,
		Description:          "test description",
		UserID:               4,
		NotificationInterval: time.Second,
	}

	ctx := context.Background()

	store.EXPECT().Create(ctx, gomock.Eq(event)).Return("", nil)

	_, err := store.Create(ctx, event)

	require.NoError(t, err)
}
