package sqlstorage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
)

const Alias = "sql"

type Config interface {
	GetStorageDSN() string
}

type Storage struct {
	Config Config
	db     *sqlx.DB
}

func New(config Config) *Storage {
	db, err := sqlx.ConnectContext(context.Background(), "postgres", config.GetStorageDSN())
	if err != nil {
		// return nil, fmt.Errorf("can not open db, %w", err)
	}

	return &Storage{
		config,
		db,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
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
