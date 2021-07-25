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

func New() *Storage {
	return &Storage{
		events: make(map[int64]storage.Event),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.increment++
	event.ID = s.increment
	s.events[event.ID] = event

	return event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {

	var err error = nil

	if _, ok := s.events[event.ID]; ok {
		s.mu.Lock()
		s.events[event.ID] = event
		s.mu.Unlock()
	} else {
		event, err = s.CreateEvent(ctx, event)
	}

	return event, err
}

func (s *Storage) RemoveEvent(ctx context.Context, event storage.Event) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; ok {
		delete(s.events, event.ID)
	}

	return nil
}

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
