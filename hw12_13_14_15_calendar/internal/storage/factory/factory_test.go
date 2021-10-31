package factorystorage

import (
	"context"
	"testing"

	internalconfig "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("storage factory", func(t *testing.T) {
		config, err := internalconfig.NewConfig("../../../configs/calendar.toml")
		if err != nil {
			t.Fatal(err)
		}

		config.Storage.Implementation = "memory"

		storage, err := GetStorage(context.Background(), config)
		if err != nil {
			t.Fatal(err)
		}

		require.IsType(t, storage, &memorystorage.Storage{}, "type of storage must be: *memorystorage.Storage")
	})
}
