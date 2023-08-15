package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/mq/rabbitmq"
	"golang.org/x/exp/slog"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/sender_config.toml", "Path to configuration file")
}

func main() {
	cfg, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("rabbit sender config error: %s", err.Error())
	}

	logg := logger.NewLogger(cfg.Logger)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	sender, err := rabbitmq.NewSender(cfg.MQ, logg)
	if err != nil {
		logg.Error("error creating rabbit sender",
			slog.String("error", err.Error()),
			slog.String("address", cfg.MQ.Host+":"+strconv.Itoa(cfg.MQ.Port)))
		cancel()
	}

	go func() {
		<-ctx.Done()

		if err := sender.Shutdown(); err != nil {
			logg.Error("error stopping rabbit sender",
				slog.String("error", err.Error()),
				slog.String("address", cfg.MQ.Host+":"+strconv.Itoa(cfg.MQ.Port)))
		}

		logg.Info("rabbit sender is stopped")
	}()

	err = sender.Consume()
	if err != nil {
		logg.Error("error consuming rabbit",
			slog.String("errors", err.Error()))
		cancel()
	}
}
