package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/mq"
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

	consumer := mq.NewConsumer(sender)

	go func() {
		<-ctx.Done()

		if err := consumer.Shutdown(); err != nil {
			logg.Error("error stopping rabbit sender",
				slog.String("error", err.Error()),
				slog.String("address", cfg.MQ.Host+":"+strconv.Itoa(cfg.MQ.Port)))
		}

		logg.Info("rabbit sender is stopped")
	}()

	notifications, err := consumer.Consume()
	if err != nil {
		logg.Error("error consuming rabbit",
			slog.String("errors", err.Error()))
		cancel()
	}

	for notification := range notifications {
		if notification.Err != nil {
			logg.Error("error receiving notification", slog.String("error", notification.Err.Error()))
			continue
		}
		logg.Info("notification is received", slog.Any("notification", notification.Message))
	}
}
