package memorystorage

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/models"
	"github.com/stretchr/testify/require"
)

func TestStorageCreate(t *testing.T) {
	st := NewStorageMemory()
	ctx := context.Background()

	events := generateEvents("test create")

	for j, event := range events {
		id, _ := st.Create(ctx, event)
		events[j].ID = id
	}

	for _, event := range events {
		if eventCreated, ok := st.events[event.ID]; !ok || eventCreated != event {
			t.Errorf("event not match")
		}
	}
}

func TestStorageUpdate(t *testing.T) {
	st := NewStorageMemory()
	ctx := context.Background()

	eventBefore := generateEvents("test update before")
	eventAfter := generateEvents("test update after")

	eventsResult := make([]models.Event, len(eventAfter))

	for j, event := range eventBefore {
		id, _ := st.Create(ctx, event)
		eventAfter[j].ID = id
	}

	for j, event := range eventAfter {
		updatedEvent, err := st.Update(ctx, event.ID, event)
		require.NoError(t, err)
		eventsResult[j] = updatedEvent
	}

	for j, updatedEvent := range eventsResult {
		if event, ok := st.events[updatedEvent.ID]; !ok || eventAfter[j] != event {
			t.Errorf("event not match")
		}
	}
}

func TestStorageUpdateError(t *testing.T) {
	st := NewStorageMemory()
	ctx := context.Background()

	eventBefore := generateEvents("test update before")
	eventAfter := generateEvents("test update after")

	for _, event := range eventBefore {
		_, _ = st.Create(ctx, event)
	}

	for _, event := range eventAfter {
		updatedEvent, err := st.Update(ctx, event.ID, event)
		require.Error(t, err)
		require.Equal(t, fmt.Errorf("updating: no event with id %s", event.ID), err)
		require.Equal(t, models.Event{}, updatedEvent)
	}
}

func TestStorageDelete(t *testing.T) {
	st := NewStorageMemory()
	ctx := context.Background()

	events := generateEvents("test delete")

	IDs := make([]string, len(events))

	for j, event := range events {
		id, _ := st.Create(ctx, event)
		IDs[j] = id
	}

	for _, id := range IDs {
		deletedID, err := st.Delete(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, deletedID)
	}

	if len(st.events) != 0 {
		t.Errorf("must be empty")
	}
}

func TestStorageDeleteError(t *testing.T) {
	st := NewStorageMemory()
	ctx := context.Background()

	events := generateEvents("test delete")

	IDs := make([]string, len(events))

	for j, event := range events {
		id, _ := st.Create(ctx, event)
		IDs[j] = id + "suffix" // create nonexistent id
	}

	for _, id := range IDs {
		deletedID, err := st.Delete(ctx, id)
		require.Error(t, err)
		require.Equal(t, fmt.Errorf("deleting: no event with id %s", id), err)
		require.Equal(t, "", deletedID)
	}

	if len(st.events) != 100 {
		t.Errorf("must be full")
	}
}

func TestStorageGetAllByDay(t *testing.T) {
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
				_, _ = st.Create(ctx, event)
			}

			actualDates := make([]time.Time, 0)

			actualEvents, err := st.GetAllByDay(ctx, tc.day)
			require.NoError(t, err)

			for _, events := range actualEvents {
				actualDates = append(actualDates, events.Date)
			}

			require.Equal(t, tc.totalDays, len(actualEvents))
			require.True(t, reflect.DeepEqual(actualDates, tc.expected))
		})
	}
}

func TestStorageGetAllByWeek(t *testing.T) {
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
				_, _ = st.Create(ctx, event)
			}

			expectedDates := tc.expected(tc.fromDay)

			actualDates := make([]time.Time, 0)

			actualEvents, err := st.GetAllByWeek(ctx, tc.fromDay)
			require.NoError(t, err)

			for _, events := range actualEvents {
				actualDates = append(actualDates, events.Date)
			}

			sort.Slice(actualDates, func(i, j int) bool {
				return actualDates[i].Before(actualDates[j])
			})

			require.True(t, reflect.DeepEqual(actualDates, expectedDates))
			require.Equal(t, tc.totalDays, len(actualEvents))
		})
	}
}

func TestStorageGetAllByMonth(t *testing.T) {
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
				_, _ = st.Create(ctx, event)
			}

			expectedDates := tc.expected(tc.fromDay)

			actualDates := make([]time.Time, 0)

			actualEvents, err := st.GetAllByMonth(ctx, tc.fromDay)
			require.NoError(t, err)

			for _, events := range actualEvents {
				actualDates = append(actualDates, events.Date)
			}

			sort.Slice(actualDates, func(i, j int) bool {
				return actualDates[i].Before(actualDates[j])
			})

			require.True(t, reflect.DeepEqual(actualDates, expectedDates))
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
