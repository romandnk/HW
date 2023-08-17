package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/mq/rabbitmq"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/service"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/postgres"
	"golang.org/x/exp/slog"
)

var configFile string

var (
	memSt      = "memory"
	postgresSt = "postgres"
)

var ErrInvalidStorageType = errors.New("invalid storage type")

func init() {
	flag.StringVar(&configFile, "config", "./configs/scheduler_config.toml", "Path to configuration file")
}

//nolint:gocognit
func main() {
	cfg, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("rabbit scheduler config error: %s", err.Error())
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg := logger.NewLogger(cfg.Logger)

	var st storage.Storage

	// use memory storage or sql storage
	switch cfg.StorageType {
	case memSt:
		st = memorystorage.NewStorageMemory()
		logg.Info("use memory scheduler storage")
	case postgresSt:
		postgresStorage := postgres.NewStoragePostgres()
		err = postgresStorage.Connect(ctx, cfg.Storage)
		if err != nil {
			logg.Error("error connecting scheduler db",
				slog.String("error", err.Error()),
				slog.String("address", cfg.Storage.Host+":"+cfg.Storage.Port))
			os.Exit(1) //nolint:gocritic
		}
		defer postgresStorage.Close()

		st = postgresStorage

		logg.Info("use postgres scheduler storage")
	default:
		logg.Error("scheduler storage", slog.String("error", ErrInvalidStorageType.Error()))
		os.Exit(1)
	}

	services := service.NewService(st)

	scheduler, err := rabbitmq.NewProducer(cfg.MQ, logg)
	if err != nil {
		cancel()
		logg.Error("error connecting scheduler rabbit",
			slog.String("errors", err.Error()),
			slog.String("address", cfg.MQ.Host+":"+strconv.Itoa(cfg.MQ.Port)))
	}

	tickerScheduler := time.NewTicker(cfg.TimeToSchedule)
	tickerDeleteOutdated := time.NewTicker(cfg.TimeToDeleteOutdated)
	done := make(chan struct{})

	go func() {
		<-ctx.Done()

		if err := scheduler.Shutdown(); err != nil {
			logg.Error("error stopping rabbit scheduler",
				slog.String("error", err.Error()),
				slog.String("address", cfg.MQ.Host+":"+strconv.Itoa(cfg.MQ.Port)))
		}

		done <- struct{}{}

		logg.Info("rabbit sender is stopped")
	}()

	for {
		select {
		case <-tickerScheduler.C:
			notifications, err := services.Notification.GetNotificationInAdvance(ctx)
			if err != nil {
				logg.Error("error getting notification", slog.String("error", err.Error()))
			}

			now := time.Now()

			for _, notification := range notifications {
				notificationMessageDate := notification.Date

				if !(notificationMessageDate.Before(now.Add(-time.Second*5)) && notificationMessageDate.After(now.Add(time.Second*5))) {
					continue
				}

				msg := rabbitmq.Message{
					EventID: notification.EventID,
					Title:   notification.Title,
					Date:    notification.Date,
					UserID:  notification.UserID,
				}

				body, err := json.Marshal(msg)
				if err != nil {
					logg.Error("error marshal notification",
						slog.Any("notification", msg),
						slog.String("error", err.Error()))
				}

				err = scheduler.Publish(ctx, body)
				if err != nil {
					logg.Error("error publish notification",
						slog.Any("notification", msg),
						slog.String("error", err.Error()))
				}
			}
		case <-tickerDeleteOutdated.C:
			err := services.Event.DeleteOutdatedEvents(ctx)
			if err != nil {
				logg.Error("error deleting outdated events", slog.String("error", err.Error()))
			}
		case <-done:
			return
		}
	}
}
