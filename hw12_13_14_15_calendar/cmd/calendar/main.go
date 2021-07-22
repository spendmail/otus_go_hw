package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/server/http"
	internalstorage "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
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

	config := internalconfig.NewConfig(configPath)
	logger := internallogger.New(config)
	storage := internalstorage.GetStorage(config)
	calendar := app.New(logger, storage)
	server := internalhttp.NewServer(config, calendar, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

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
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	logger.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
