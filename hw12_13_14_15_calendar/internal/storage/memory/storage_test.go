package memorystorage

import (
	"context"
	"testing"
	"time"

	internalstorage "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("storage memory", func(t *testing.T) {
		storage := New()

		ctx := context.Background()

		// Create event 1
		event, err := storage.CreateEvent(ctx, internalstorage.Event{
			Title:       "event 1",
			BeginDate:   time.Now().Add(1 * time.Hour),
			EndDate:     time.Now().Add(1 * time.Hour),
			Description: "description",
			OwnerID:     1,
		})

		require.NoError(t, err)
		require.IsType(t, event, internalstorage.Event{}, "type of event must be: storage.Event")

		// Create event 2
		event, err = storage.CreateEvent(ctx, internalstorage.Event{
			Title:       "event 2",
			BeginDate:   time.Now().Add(24 * 3 * time.Hour),
			EndDate:     time.Now().Add(24 * 3 * time.Hour),
			Description: "description",
			OwnerID:     1,
		})
		require.NoError(t, err)
		require.IsType(t, event, internalstorage.Event{}, "type of event must be: storage.Event")

		// Create event 3
		event, err = storage.CreateEvent(ctx, internalstorage.Event{
			Title:       "event 1",
			BeginDate:   time.Now().Add(24 * 14 * time.Hour),
			EndDate:     time.Now().Add(24 * 14 * time.Hour),
			Description: "description",
			OwnerID:     1,
		})
		require.NoError(t, err)
		require.IsType(t, event, internalstorage.Event{}, "type of event must be: storage.Event")

		// Update event 3
		title := "event 3 updated"
		event.Title = title
		event, err = storage.UpdateEvent(ctx, event)
		require.NoError(t, err)
		require.IsType(t, event, internalstorage.Event{}, "type of event must be: storage.Event")
		require.Equal(t, event.Title, title, "title has not been updated")

		// Getting 1 event
		events, err := storage.GetDayAheadEvents(ctx)
		require.NoError(t, err)
		require.Len(t, events, 1, "Len is %s, but expected %s", len(events), 1)

		// Getting 2 events
		events, err = storage.GetWeekAheadEvents(ctx)
		require.NoError(t, err)
		require.Len(t, events, 2, "Len is %s, but expected %s", len(events), 2)

		// Getting 3 events
		events, err = storage.GetMonthAheadEvents(ctx)
		require.NoError(t, err)
		require.Len(t, events, 3, "Len is %s, but expected %s", len(events), 3)

		// Remove event
		err = storage.RemoveEvent(ctx, event)
		require.NoError(t, err)

		// Getting 2 events
		events, err = storage.GetMonthAheadEvents(ctx)
		require.NoError(t, err)
		require.Len(t, events, 2, "Len is %s, but expected %s", len(events), 2)
	})
}
