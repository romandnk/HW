package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/grpc"
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
		log.Fatalf("config errors: %s", err.Error())
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
			logg.Error("errors connecting db",
				slog.String("errors", err.Error()),
				slog.String("address", config.Storage.DB.Host+":"+config.Storage.DB.Port))
			cancel()
		}
		defer db.Close()

		st = sqlstorage.NewStorageSQL(db)
	}

	services := service.NewService(st)

	handlerHTTP := internalhttp.NewHandlerHTTP(services, logg)
	handlerGRPC := grpc.NewHandlerGRPC(services, logg)

	serverHTTP := internalhttp.NewServerHTTP(config.ServerHTTP, handlerHTTP.InitRoutes())
	serverGRPC := grpc.NewServerGRPC(handlerGRPC, logg, config.ServerGRPC)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverHTTP.Stop(ctx); err != nil {
			logg.Error("errors stopping HTTPServer",
				slog.String("address http", net.JoinHostPort(config.ServerHTTP.Host, config.ServerHTTP.Port)))
			cancel()
		}

		serverGRPC.Stop()

		logg.Info("calendar is stopped")
	}()

	logg.Info("calendar is running...",
		slog.String("address http", net.JoinHostPort(config.ServerHTTP.Host, config.ServerHTTP.Port)),
		slog.String("address grpc", net.JoinHostPort(config.ServerGRPC.Host, config.ServerGRPC.Port)))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := serverHTTP.Start(); err != nil {
			logg.Error("errors starting HTTPServer",
				slog.String("address http", net.JoinHostPort(config.ServerHTTP.Host, config.ServerHTTP.Port)))
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if err := serverGRPC.Start(config.ServerGRPC); err != nil {
			logg.Error("errors starting GRPCServer",
				slog.String("address grpc", net.JoinHostPort(config.ServerGRPC.Host, config.ServerGRPC.Port)))
			cancel()
		}
	}()

	wg.Wait()
}
