package factorystorage

import (
	"context"

	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type Config interface {
	GetStorageImplementation() string
	GetStorageDSN() string
}

func GetStorage(ctx context.Context, config Config) (app.Storage, error) {
	switch config.GetStorageImplementation() {
	case sqlstorage.Alias:
		storage := sqlstorage.New(config)
		err := storage.Connect(ctx)

		return storage, err

	case memorystorage.Alias:
		return memorystorage.New(), nil
	default:
		return memorystorage.New(), nil
	}
}
