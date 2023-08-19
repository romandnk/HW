package postgres

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"
)

func TestStorageUpdateScheduledNotification(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	id := uuid.New().String()

	ctx := context.Background()

	storage := NewStoragePostgres()
	storage.db = mock

	query := fmt.Sprintf(`
		UPDATE %s
		SET scheduled = true 
		WHERE id = $1`, eventsTable)

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = storage.UpdateScheduledNotification(ctx, id)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestStorageUpdateScheduledNotificationError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	id := uuid.New().String()

	ctx := context.Background()

	storage := NewStoragePostgres()
	storage.db = mock

	query := fmt.Sprintf(`
		UPDATE %s
		SET scheduled = true 
		WHERE id = $1`, eventsTable)

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err = storage.UpdateScheduledNotification(ctx, id)
	expectedError := "notification wasn't updated with id: " + id
	require.EqualError(t, err, expectedError)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}
