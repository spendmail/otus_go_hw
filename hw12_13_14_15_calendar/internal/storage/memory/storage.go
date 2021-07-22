package memorystorage

import (
	"sync"
	"time"

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

func (s *Storage) AddEvent(event storage.Event) error {
	return nil
}

func (s *Storage) UpdateEvent(event storage.Event) error {
	return nil
}

func (s *Storage) RemoveEvent(eventID string) error {
	return nil
}

func (s *Storage) DailyEvents(date time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}

func (s *Storage) WeeklyEvents(date time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}

func (s *Storage) MonthEvents(date time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}
