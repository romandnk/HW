package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
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

	_ = service.NewService(st)

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

	for {
		select {
		case <-ctx.Done():
			if err := scheduler.CloseChannel(); err != nil {
				logg.Error("error closing rabbit scheduler channel", slog.String("errors", err.Error()))
			}

			if err := scheduler.CloseConn(); err != nil {
				logg.Error("error closing rabbit scheduler connection",
					slog.String("errors", err.Error()),
					slog.String("address", cfg.MQ.Host+":"+strconv.Itoa(cfg.MQ.Port)))
			}

			logg.Info("rabbit scheduler is stopped")
			return
		default:

			// TODO: notification logic

			time.Sleep(time.Second * 2)
			fmt.Println("1")
		}
	}
}
