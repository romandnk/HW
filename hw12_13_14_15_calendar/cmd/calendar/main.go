package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/http"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/service"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/sql"
	"golang.org/x/exp/slog"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/calendar_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("config error: %s", err.Error())
	}

	logg := logger.NewLogger(config.Logger.Level, config.Logger.Representation)

	logg.Info("use logging")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var st storage.Storage

	// use memory storage or sql storage
	switch config.Storage.Type {
	case "memory":
		st = memorystorage.NewStorageMemory()
	case "postgres":
		db, err := sqlstorage.NewPostgresDB(ctx, config.Storage.DB)
		if err != nil {
			logg.Error("error connecting db",
				slog.String("error", err.Error()),
				slog.String("address", config.Storage.DB.Host+":"+config.Storage.DB.Port))
			os.Exit(1) //nolint:gocritic
		}
		defer db.Close()

		st = sqlstorage.NewStorageSQL(db)
	}

	services := service.NewService(st)

	handler := internalhttp.NewHandler(services, logg)

	server := internalhttp.NewServer(config.Server, handler.InitRoutes())

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("error stopping server", slog.String("address", net.JoinHostPort(config.Server.Host, config.Server.Port)))
			cancel()
			os.Exit(1)
		}

		logg.Info("calendar is stopped")
	}()

	logg.Info("calendar is running...", slog.String("address", net.JoinHostPort(config.Server.Host, config.Server.Port)))

	if err := server.Start(); err != nil {
		logg.Error("error starting server", slog.String("address", net.JoinHostPort(config.Server.Host, config.Server.Port)))
		cancel()
	}
}
