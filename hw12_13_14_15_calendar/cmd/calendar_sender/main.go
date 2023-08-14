package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/cmd/config"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/mq/rabbitmq"
	"golang.org/x/exp/slog"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/sender_config.toml", "Path to configuration file")
}

func main() {
	cfg, err := config.NewSenderConfig(configFile)
	if err != nil {
		log.Fatalf("rabbit sender config error: %s", err.Error())
	}

	logg := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Representation)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	sender, err := rabbitmq.NewSender(cfg.MQ, logg)
	if err != nil {
		cancel()
		logg.Error("error connecting sender rabbit",
			slog.String("errors", err.Error()),
			slog.String("address", cfg.MQ.Host+":"+strconv.Itoa(cfg.MQ.Port)))
	}

	err = sender.OpenChannel()
	if err != nil {
		cancel()
		logg.Error("error opening channel sender rabbit",
			slog.String("errors", err.Error()))
	}

	go func() {
		<-ctx.Done()

		if err := sender.CloseConn(); err != nil {
			logg.Error("error closing rabbit sender connection",
				slog.String("errors", err.Error()),
				slog.String("address", cfg.MQ.Host+":"+strconv.Itoa(cfg.MQ.Port)))
		}

		if err := sender.CloseChannel(); err != nil {
			logg.Error("error closing rabbit sender channel", slog.String("errors", err.Error()))
		}

		logg.Info("rabbit sender is stopped")
		os.Exit(0)
	}()

	err = sender.Consume()
	if err != nil {
		cancel()
		logg.Error("error consuming rabbit",
			slog.String("errors", err.Error()))
	}

	for notification := range sender.Handle(ctx) {
		if notification.Err != nil {
			logg.Error("error receiving notification", slog.String("error", notification.Err.Error()))
			continue
		}
		logg.Info("received notification",
			slog.Any("notification", notification.Message))
	}
}
