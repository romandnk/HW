package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"github.com/stretchr/testify/require"
)

func TestStorage_GetNotificationsInAdvance(t *testing.T) {
	st := NewStorageMemory()
	ctx := context.Background()

	st.events["id1"] = models.Event{
		ID:                   "id1",
		Date:                 time.Now().Add(time.Hour),
		NotificationInterval: 0,
		Scheduled:            false,
	}
	st.events["id2"] = models.Event{
		ID:                   "id2",
		Date:                 time.Now().Add(time.Hour),
		NotificationInterval: 0,
		Scheduled:            true,
	}
	st.events["id3"] = models.Event{
		ID:                   "id3",
		Date:                 time.Now().Add(time.Minute * 30),
		NotificationInterval: 0,
		Scheduled:            false,
	}
	st.events["id4"] = models.Event{
		ID:                   "id4",
		Date:                 time.Now().Add(-time.Hour),
		NotificationInterval: 0,
		Scheduled:            false,
	}
	st.events["id5"] = models.Event{
		ID:                   "id5",
		Date:                 time.Now().Add(time.Second * 3),
		NotificationInterval: 0,
		Scheduled:            false,
	}

	notifications, err := st.GetNotificationsInAdvance(ctx)
	require.NoError(t, err)
	require.Len(t, notifications, 3)

	require.Equal(t, notifications[0].Date, st.events["id5"].Date)
	require.Equal(t, notifications[1].Date, st.events["id3"].Date)
	require.Equal(t, notifications[2].Date, st.events["id1"].Date)
}
