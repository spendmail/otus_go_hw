package memorystorage

import (
	"context"
	"sync"

	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
)

const Alias = "memory"

type Storage struct {
	// TODO
	mu sync.RWMutex
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	return storage.Event{}, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	return storage.Event{}, nil
}

func (s *Storage) RemoveEvent(ctx context.Context, event storage.Event) error {
	return nil
}

func (s *Storage) GetDayAheadEvents(ctx context.Context) ([]storage.Event, error) {
	return []storage.Event{}, nil
}

func (s *Storage) GetWeekAheadEvents(ctx context.Context) ([]storage.Event, error) {
	return []storage.Event{}, nil
}

func (s *Storage) GetMonthAheadEvents(ctx context.Context) ([]storage.Event, error) {
	return []storage.Event{}, nil
}
