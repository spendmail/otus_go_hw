package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
)

const Alias = "memory"

type Storage struct {
	mu        sync.RWMutex
	increment int64
	events    map[int64]storage.Event
}

// New returns a new memory storage instance.
func New() *Storage {
	return &Storage{
		events: make(map[int64]storage.Event),
	}
}

// CreateEvent saves event into a memory storage.
func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.increment++
	event.ID = s.increment
	s.events[event.ID] = event

	return event, nil
}

// UpdateEvent updates event in memory storage if exists.
func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	var err error = nil

	s.mu.RLock()
	_, isSet := s.events[event.ID] //nolint:ifshort
	s.mu.RUnlock()

	if isSet {
		s.mu.Lock()
		s.events[event.ID] = event
		s.mu.Unlock()
	} else {
		event, err = s.CreateEvent(ctx, event)
	}

	return event, err
}

// RemoveEvent removes event from memory storage if exists.
func (s *Storage) RemoveEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, event.ID)

	return nil
}

// GetDayAheadEvents returns a day events slice.
func (s *Storage) GetDayAheadEvents(ctx context.Context) ([]storage.Event, error) {
	start := time.Now()
	end := time.Now().Add(24 * time.Hour)
	var events []storage.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		if event.BeginDate.After(start) && event.BeginDate.Before(end) {
			events = append(events, event)
		}
	}

	return events, nil
}

// GetWeekAheadEvents returns a week events slice.
func (s *Storage) GetWeekAheadEvents(ctx context.Context) ([]storage.Event, error) {
	start := time.Now()
	end := time.Now().Add(24 * 7 * time.Hour)
	var events []storage.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		if event.BeginDate.After(start) && event.BeginDate.Before(end) {
			events = append(events, event)
		}
	}

	return events, nil
}

// GetMonthAheadEvents returns a month events slice.
func (s *Storage) GetMonthAheadEvents(ctx context.Context) ([]storage.Event, error) {
	start := time.Now()
	end := time.Now().Add(24 * 7 * 30 * time.Hour)
	var events []storage.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		if event.BeginDate.After(start) && event.BeginDate.Before(end) {
			events = append(events, event)
		}
	}

	return events, nil
}
