package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	internalconfig "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/server/grpc"
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
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	// Application initialization.
	calendar := app.New(logger, storage)

	// HTTP Server initialization.
	httpServer := internalhttp.NewServer(config, calendar, logger)

	// GRPC Server initialization.
	grpcServer := internalgrpc.NewServer(config, calendar, logger)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		signalNotifyCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP)
		defer stop()

		// Locking until OS signal is sent or context cancel func is called.
		select {
		case <-ctx.Done():
			return
		case <-signalNotifyCtx.Done():
		}

		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			logger.Error(err.Error())
		}

		grpcServer.Stop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Info("starting grpc server...")

		// Locking over here until server is stopped.
		if err := grpcServer.Start(); err != nil {
			logger.Error(err.Error())
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Info("starting http server...")

		// Locking over here until server is stopped.
		if err := httpServer.Start(); err != nil {
			logger.Error(err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	logger.Info("calendar is running...")

	wg.Wait()
}
