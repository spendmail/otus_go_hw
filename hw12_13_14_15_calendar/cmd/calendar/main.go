package main

import (
	"context"
	"flag"
	"log"
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

	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()

	// Storage initialization.
	storage, err := factorystorage.GetStorage(mainCtx, config)
	if err != nil {
		logger.Error(err.Error())
		mainCancel()
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

		signalNotifyCtx, signalNotifyStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP)
		defer signalNotifyStop()

		// Locking until OS signal is sent or context mainCancel func is called.
		select {
		case <-mainCtx.Done():
			return
		case <-signalNotifyCtx.Done():
		}

		mainCancel()

		// Stopping http server.
		stopHTTPCtx, stopHTTPCancel := context.WithTimeout(context.Background(), time.Second*3)
		defer stopHTTPCancel()
		if err := httpServer.Stop(stopHTTPCtx); err != nil {
			logger.Error(err.Error())
		}

		// Stopping grpc server.
		grpcServer.Stop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Info("starting grpc server...")

		// Locking over here until server is stopped.
		if err := grpcServer.Start(); err != nil {
			logger.Error(err.Error())
			mainCancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Info("starting http server...")

		// Locking over here until server is stopped.
		if err := httpServer.Start(); err != nil {
			logger.Error(err.Error())
			mainCancel()
		}
	}()

	logger.Info("calendar is running...")

	wg.Wait()
}
