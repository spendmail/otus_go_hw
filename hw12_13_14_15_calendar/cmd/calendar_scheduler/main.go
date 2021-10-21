package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	internalconfig "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	internalrabbitmq "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/rabbitmq"
	factorystorage "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage/factory"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "/etc/calendar_scheduler/calendar_scheduler.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP)
	defer cancel()

	// Config initialization.
	config, err := internalconfig.NewConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Logger initialization.
	logger := internallogger.New(config)

	// Storage initialization.
	storage, err := factorystorage.GetStorage(ctx, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// RabbitMQ client initialization.
	rabbitClient, err := internalrabbitmq.NewClient(config, logger)
	if err != nil {
		fmt.Println(err)
		return
	}

	// RabbitMQ exchange initialization.
	err = rabbitClient.DeclareExchange()
	if err != nil {
		fmt.Println(err)
		return
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			// Getting suitable events to be notified.
			events, err := storage.GetComingEvents(ctx)
			if err != nil {
				logger.Error(err.Error())
			}

			if len(events) == 0 {
				continue
			}

			for _, event := range events {
				// Sending event into a queue broker.
				err = rabbitClient.SendEventNotification(event)
				if err != nil {
					logger.Error(err.Error())
				} else {
					// If notification has been successfully sent, setting NotificationSent flag.
					event.NotificationSent = true
					event, err = storage.UpdateEvent(ctx, event)
					if err != nil {
						logger.Error(err.Error())
					}
				}

				// Removing expired events.
				err = storage.RemoveExpiredEvents(ctx)
				if err != nil {
					logger.Error(err.Error())
				}
			}
		}
	}()

	logger.Info("calendar scheduler is running...")

	<-ctx.Done()
}
