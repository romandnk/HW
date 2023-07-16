package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/http"
	st "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/sql"
	"golang.org/x/exp/slog"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./config/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)

	log := logger.NewLogger(config.Logger.Level, config.Logger.Representation)

	log.Info("use logging")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var storage st.StoreEvent

	// use memory storage or sql storage
	if config.Storage.Memory {
		storage = memorystorage.NewStorageMemory()
	} else {
		db, err := sqlstorage.NewPostgresDB(ctx, config.Storage.DB)
		if err != nil {
			log.Error("error connecting db", slog.String("address", config.Storage.DB.Host+":"+config.Storage.DB.Port))
			cancel()
			os.Exit(1)
		}
		defer db.Close()

		storage = sqlstorage.NewStorageSQL(db)
	}

	handler := internalhttp.NewHandler(storage)

	server := internalhttp.NewServer(config.Server, handler.InitRoutes(log))

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error("error stopping server", slog.String("address", net.JoinHostPort(config.Server.Host, config.Server.Port)))
			cancel()
			os.Exit(1)
		}

		log.Info("calendar is stopped")
	}()

	log.Info("calendar is running...", slog.String("address", net.JoinHostPort(config.Server.Host, config.Server.Port)))

	if err := server.Start(); err != nil {
		log.Error("error starting server", slog.String("address", net.JoinHostPort(config.Server.Host, config.Server.Port)))
		cancel()
		os.Exit(1)
	}
}
