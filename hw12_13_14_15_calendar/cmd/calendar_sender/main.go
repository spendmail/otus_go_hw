package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/stdlib"
	internalconfig "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	internalrabbitmq "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/rabbitmq"
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

	// RabbitMQ queue initialization.
	queue, err := rabbitClient.DeclareQueue()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Binding queue with appropriate exchange.
	err = rabbitClient.BindQueue(queue)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Getting notification channel.
	messages, err := rabbitClient.Consume(queue)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Getting notifications.
	go func() {
		for d := range messages {
			notification := internalrabbitmq.Notification{}
			err := json.Unmarshal(d.Body, &notification)
			if err != nil {
				logger.Error(err.Error())
			} else {
				// Fake sending for receiving notification.
				SendNotification(notification)
			}
		}
	}()

	logger.Info("calendar sender is running...")

	<-ctx.Done()
}

// SendNotification sends notification to a fake recipient.
func SendNotification(notification internalrabbitmq.Notification) {
	fmt.Printf("Notification %v has been successfully sent\n", notification)
}
