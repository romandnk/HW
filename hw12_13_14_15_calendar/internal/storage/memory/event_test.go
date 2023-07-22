package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"github.com/stretchr/testify/require"
)

func TestStorageCreateEvent(t *testing.T) {
	st := NewStorageMemory()
	ctx := context.Background()

	events := generateEvents("test create")

	for j, event := range events {
		id, err := st.CreateEvent(ctx, event)
		require.NoError(t, err)
		events[j].ID = id
	}

	for _, event := range events {
		if eventCreated, ok := st.events[event.ID]; !ok || eventCreated != event {
			t.Errorf("event not match")
		}
	}
}

func TestStorageUpdateEvent(t *testing.T) {
	st := NewStorageMemory()
	ctx := context.Background()

	eventBefore := generateEvents("test update before")
	eventAfter := generateEvents("test update after")

	eventsResult := make([]models.Event, len(eventAfter))

	for j, event := range eventBefore {
		id, err := st.CreateEvent(ctx, event)
		require.NoError(t, err)
		eventAfter[j].ID = id
	}

	for j, event := range eventAfter {
		updatedEvent, err := st.UpdateEvent(ctx, event.ID, event)
		require.NoError(t, err)
		eventsResult[j] = updatedEvent
	}

	for j, updatedEvent := range eventsResult {
		if event, ok := st.events[updatedEvent.ID]; !ok || eventAfter[j] != event {
			t.Errorf("event not match")
		}
	}
}

func TestStorageUpdateEventError(t *testing.T) {
	st := NewStorageMemory()
	ctx := context.Background()

	eventBefore := generateEvents("test update before")
	eventAfter := generateEvents("test update after")

	for _, event := range eventBefore {
		_, err := st.CreateEvent(ctx, event)
		require.NoError(t, err)
	}

	for _, event := range eventAfter {
		updatedEvent, err := st.UpdateEvent(ctx, event.ID, event)
		require.Error(t, err)
		require.EqualError(t, err, fmt.Errorf("no event with id %s", event.ID).Error())
		require.Equal(t, models.Event{}, updatedEvent)
	}
}

func TestStorageDeleteEvent(t *testing.T) {
	st := NewStorageMemory()
	ctx := context.Background()

	events := generateEvents("test delete")

	IDs := make([]string, len(events))

	for j, event := range events {
		id, err := st.CreateEvent(ctx, event)
		require.NoError(t, err)
		IDs[j] = id
	}

	for _, id := range IDs {
		err := st.DeleteEvent(ctx, id)
		require.NoError(t, err)
	}

	require.Len(t, st.events, 0, "must be empty")
}

func TestStorageDeleteEventError(t *testing.T) {
	st := NewStorageMemory()
	ctx := context.Background()

	events := generateEvents("test delete")

	IDs := make([]string, len(events))

	for j, event := range events {
		id, err := st.CreateEvent(ctx, event)
		require.NoError(t, err)
		IDs[j] = id + "suffix" // create nonexistent id
	}

	for _, id := range IDs {
		err := st.DeleteEvent(ctx, id)
		require.Error(t, err)
		require.EqualError(t, err, fmt.Errorf("no event with id %s", id).Error())
	}

	require.Len(t, st.events, 100, "must be full")
}

func TestStorageGetAllByDayEvents(t *testing.T) {
	testCases := []struct {
		day       time.Time
		expected  []time.Time
		totalDays int
	}{
		{
			day:       time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local),
			expected:  []time.Time{time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)},
			totalDays: 1,
		},
		{
			day:       time.Date(1999, 12, 31, 0, 0, 0, 0, time.Local),
			expected:  []time.Time{},
			totalDays: 0,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("test get all by day", func(t *testing.T) {
			st := NewStorageMemory()
			ctx := context.Background()

			events := generateEvents("test get all by day")

			for _, event := range events {
				_, _ = st.CreateEvent(ctx, event)
			}

			actualDates := make([]time.Time, 0)

			actualEvents, err := st.GetAllByDayEvents(ctx, tc.day)
			require.NoError(t, err)

			for _, events := range actualEvents {
				actualDates = append(actualDates, events.Date)
			}

			require.Equal(t, tc.totalDays, len(actualEvents))
			require.Equal(t, tc.expected, actualDates)
		})
	}
}

func TestStorageGetAllByWeekEvents(t *testing.T) {
	testCases := []struct {
		fromDay   time.Time
		expected  func(day time.Time) []time.Time
		totalDays int
	}{
		{
			fromDay: time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local),
			expected: func(day time.Time) []time.Time {
				var eventsDate []time.Time
				for i := 0; i < 7; i++ {
					eventsDate = append(eventsDate, day)
					day = day.AddDate(0, 0, 1)
				}
				return eventsDate
			},
			totalDays: 7,
		},
		{
			fromDay: time.Date(1999, 12, 29, 0, 0, 0, 0, time.Local),
			expected: func(day time.Time) []time.Time {
				var eventsDate []time.Time
				day = day.AddDate(0, 0, 4)
				for i := 0; i < 3; i++ {
					eventsDate = append(eventsDate, day)
					day = day.AddDate(0, 0, 1)
				}
				return eventsDate
			},
			totalDays: 3,
		},
		{
			fromDay: time.Date(2000, 4, 6, 0, 0, 0, 0, time.Local),
			expected: func(day time.Time) []time.Time {
				var eventsDate []time.Time
				for i := 0; i < 5; i++ {
					eventsDate = append(eventsDate, day)
					day = day.AddDate(0, 0, 1)
				}
				return eventsDate
			},
			totalDays: 5,
		},
		{
			fromDay: time.Date(1999, 4, 6, 0, 0, 0, 0, time.Local),
			expected: func(day time.Time) []time.Time {
				eventsDate := make([]time.Time, 0)
				return eventsDate
			},
			totalDays: 0,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("get all by week", func(t *testing.T) {
			st := NewStorageMemory()
			ctx := context.Background()

			events := generateEvents("test get all by week")

			for _, event := range events {
				_, _ = st.CreateEvent(ctx, event)
			}

			expectedDates := tc.expected(tc.fromDay)

			actualDates := make([]time.Time, 0)

			actualEvents, err := st.GetAllByWeekEvents(ctx, tc.fromDay)
			require.NoError(t, err)

			for _, events := range actualEvents {
				actualDates = append(actualDates, events.Date)
			}

			require.ElementsMatch(t, expectedDates, actualDates)
			require.Equal(t, tc.totalDays, len(actualEvents))
		})
	}
}

func TestStorageGetAllByMonthEvents(t *testing.T) {
	testCases := []struct {
		fromDay   time.Time
		expected  func(day time.Time) []time.Time
		totalDays int
	}{
		{
			fromDay: time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local),
			expected: func(day time.Time) []time.Time {
				eventsDate := make([]time.Time, 0, 32)
				for i := 0; i < 30; i++ {
					eventsDate = append(eventsDate, day)
					day = day.AddDate(0, 0, 1)
				}
				return eventsDate
			},
			totalDays: 30,
		},
		{
			fromDay: time.Date(1999, 12, 29, 0, 0, 0, 0, time.Local),
			expected: func(day time.Time) []time.Time {
				eventsDate := make([]time.Time, 0, 32)
				day = day.AddDate(0, 0, 4)
				for i := 0; i < 26; i++ {
					eventsDate = append(eventsDate, day)
					day = day.AddDate(0, 0, 1)
				}
				return eventsDate
			},
			totalDays: 26,
		},
		{
			fromDay: time.Date(2000, 4, 6, 0, 0, 0, 0, time.Local),
			expected: func(day time.Time) []time.Time {
				eventsDate := make([]time.Time, 0, 32)
				for i := 0; i < 5; i++ {
					eventsDate = append(eventsDate, day)
					day = day.AddDate(0, 0, 1)
				}
				return eventsDate
			},
			totalDays: 5,
		},
		{
			fromDay: time.Date(1999, 4, 6, 0, 0, 0, 0, time.Local),
			expected: func(day time.Time) []time.Time {
				eventsDate := make([]time.Time, 0)
				return eventsDate
			},
			totalDays: 0,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("get all by month", func(t *testing.T) {
			st := NewStorageMemory()
			ctx := context.Background()

			events := generateEvents("test get all by month")

			for _, event := range events {
				_, _ = st.CreateEvent(ctx, event)
			}

			expectedDates := tc.expected(tc.fromDay)

			actualDates := make([]time.Time, 0)

			actualEvents, err := st.GetAllByMonthEvents(ctx, tc.fromDay)
			require.NoError(t, err)

			for _, events := range actualEvents {
				actualDates = append(actualDates, events.Date)
			}

			require.ElementsMatch(t, expectedDates, actualDates)
			require.Equal(t, tc.totalDays, len(actualEvents))
		})
	}
}

func generateEvents(titleText string) []models.Event {
	var events []models.Event

	currentDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)

	for i := 1; i <= 100; i++ {
		currentDate = currentDate.AddDate(0, 0, 1)

		event := models.Event{
			ID:                   uuid.New().String(),
			Title:                fmt.Sprintf("%s %d", titleText, i),
			Date:                 currentDate,
			Duration:             time.Duration(i),
			Description:          "",
			UserID:               i,
			NotificationInterval: time.Duration(i),
		}

		events = append(events, event)
	}

	return events
}
