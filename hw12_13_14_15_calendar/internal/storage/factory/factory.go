package factorystorage

import (
	"time"

	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type Storage interface {
	AddEvent(event storage.Event) error
	UpdateEvent(event storage.Event) error
	RemoveEvent(eventID string) error
	DailyEvents(date time.Time) ([]storage.Event, error)
	WeeklyEvents(date time.Time) ([]storage.Event, error)
	MonthEvents(date time.Time) ([]storage.Event, error)
}

type Config interface {
	GetStorageImplementation() string
	GetStorageDSN() string
}

func GetStorage(config Config) Storage {
	var implementation Storage

	switch config.GetStorageImplementation() {
	case sqlstorage.Alias:
		implementation = sqlstorage.New(config)
	case memorystorage.Alias:
		implementation = memorystorage.New()
	default:
		implementation = memorystorage.New()
	}

	return implementation
}
