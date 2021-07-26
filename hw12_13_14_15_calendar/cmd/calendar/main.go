package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/server/http"
	factorystorage "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage/factory"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Config initialization.
	config, err := internalconfig.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Logger initialization.
	logger := internallogger.New(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Storage initialization.
	storage, err := factorystorage.GetStorage(ctx, config)
	if err != nil {
		logger.Error(err.Error())
	}

	// Application initialization.
	calendar := app.New(logger, storage)

	// HTTP server initialization.
	server := internalhttp.NewServer(config, calendar, logger)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

		// Locking until OS signal is sent or context cancel func is called.
		select {
		case <-ctx.Done():
			return
		case <-signals:
		}

		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Error(err.Error())
		}
	}()

	logger.Info("calendar is running...")

	// Locking till server is listening the socket.
	if err := server.Start(ctx); err != nil {
		logger.Error(err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
