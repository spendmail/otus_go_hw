package storage

import (
	"time"

	memorystorage "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type Storage interface {
	AddEvent(event Event) error
	UpdateEvent(event Event) error
	RemoveEvent(eventID string) error
	DailyEvents(date time.Time) ([]Event, error)
	WeeklyEvents(date time.Time) ([]Event, error)
	MonthEvents(date time.Time) ([]Event, error)
}

type Config interface {
	GetStorageImplementation() string
	GetStorageDSN() string
}

func GetStorage(config Config) Storage {
	var storage Storage

	switch config.GetStorageImplementation() {
	case sqlstorage.Alias:
		storage = sqlstorage.New(config)
	case memorystorage.Alias:
		storage = memorystorage.New()
	default:
		storage = memorystorage.New()
	}

	return storage
}
