package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	internalconfig "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	internalrabbitmq "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/rabbitmq"
	"log"
	"os/signal"
	"syscall"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "/etc/calendar_sender/calendar_sender.toml", "Path to configuration file")
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP)
	defer cancel()

	rabbitClient, err := internalrabbitmq.NewClient(config, logger)
	if err != nil {
		log.Fatal(err)
	}

	err = rabbitClient.DeclareExchange()
	if err != nil {
		log.Fatal(err)
	}

	queue, err := rabbitClient.DeclareQueue()
	if err != nil {
		log.Fatal(err)
	}

	err = rabbitClient.BindQueue(queue)
	if err != nil {
		log.Fatal(err)
	}

	messages, err := rabbitClient.Consume(queue)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for d := range messages {
			notification := internalrabbitmq.Notification{}
			err := json.Unmarshal(d.Body, &notification)
			if err != nil {
				logger.Error(err.Error())
			} else {
				SendNotification(notification)
			}
		}
	}()

	logger.Info("calendar sender is running...")

	<-ctx.Done()
}

func SendNotification(notification internalrabbitmq.Notification) {
	fmt.Printf("Notification %v has been successfully sent\n", notification)
}
