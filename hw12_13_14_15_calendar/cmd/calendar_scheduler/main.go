package main

import (
	"context"
	"encoding/json"
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
	flag.StringVar(&configFile, "config", "./configs/scheduler_config.toml", "Path to configuration file")
}

func main() {
	cfg, err := config.NewSchedulerConfig(configFile)
	if err != nil {
		log.Fatalf("rabbit scheduler config error: %s", err.Error())
	}

	logg := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Representation)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	scheduler, err := rabbitmq.NewScheduler(cfg.MQ, logg)
	if err != nil {
		cancel()
		logg.Error("error connecting scheduler rabbit",
			slog.String("errors", err.Error()),
			slog.String("address", cfg.MQ.Host+":"+strconv.Itoa(cfg.MQ.Port)))
	}

	err = scheduler.OpenChannel()
	if err != nil {
		cancel()
		logg.Error("error opening channel scheduler rabbit",
			slog.String("errors", err.Error()))
	}

	go func() {
		<-ctx.Done()

		if err := scheduler.CloseConn(); err != nil {
			logg.Error("error closing rabbit scheduler connection",
				slog.String("errors", err.Error()),
				slog.String("address", cfg.MQ.Host+":"+strconv.Itoa(cfg.MQ.Port)))
		}

		if err := scheduler.CloseChannel(); err != nil {
			logg.Error("error closing rabbit scheduler channel", slog.String("errors", err.Error()))
		}

		logg.Info("rabbit scheduler is stopped")
		os.Exit(0)
	}()

	str := rabbitmq.Message{
		EventID:     "test",
		Title:       "test",
		Description: "test",
		Date:        0,
		UserID:      "test",
	}

	body, _ := json.Marshal(str)

	err = scheduler.Publish(ctx, body)
	if err != nil {
		cancel()
		logg.Error("error publishing rabbit",
			slog.String("errors", err.Error()))
	}
}
