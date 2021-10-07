package main

import (
	"context"
	"flag"
	_ "github.com/jackc/pgx/stdlib"
	internalconfig "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
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

	logger.Info("sender is running...")

	<-ctx.Done()
}
