package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"os"
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
	flag.StringVar(&configFile, "config", "./configs/calendar_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := NewConfig(configFile)
	if err != nil {
		log.Fatalf("calendar config error: %s", err.Error())
	}

	logg := logger.NewLogger(cfg.Logger)

	logg.Info("use logging")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var st storage.Storage

	// use memory storage or sql storage
	switch cfg.StorageType {
	case memSt:
		st = memorystorage.NewStorageMemory()
		logg.Info("use memory calendar storage")
	case postgresSt:
		postgresStorage := postgres.NewStoragePostgres()
		err = postgresStorage.Connect(ctx, cfg.Storage)
		if err != nil {
			logg.Error("error connecting calendar db",
				slog.String("error", err.Error()),
				slog.String("address", cfg.Storage.Host+":"+cfg.Storage.Port))
			os.Exit(1) //nolint:gocritic
		}
		defer postgresStorage.Close()

		st = postgresStorage

		logg.Info("use postgres calendar storage")
	default:
		logg.Error("calendar storage", slog.String("error", ErrInvalidStorageType.Error()))
		os.Exit(1)
	}

	services := service.NewService(st)

	handlerHTTP := internalhttp.NewHandlerHTTP(services, logg)
	handlerGRPC := grpc.NewHandlerGRPC(services, logg)

	serverHTTP := internalhttp.NewServerHTTP(cfg.ServerHTTP, handlerHTTP.InitRoutes(cfg.Logger.LogFilePath))
	serverGRPC := grpc.NewServerGRPC(handlerGRPC, logg, cfg.ServerGRPC, cfg.Logger.LogFilePath)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverHTTP.Stop(ctx); err != nil {
			logg.Error("error stopping HTTPServer",
				slog.String("address http", net.JoinHostPort(cfg.ServerHTTP.Host, cfg.ServerHTTP.Port)))
			cancel()
		}

		serverGRPC.Stop()

		logg.Info("calendar is stopped")
	}()

	logg.Info("calendar is running...",
		slog.String("address http", net.JoinHostPort(cfg.ServerHTTP.Host, cfg.ServerHTTP.Port)),
		slog.String("address grpc", net.JoinHostPort(cfg.ServerGRPC.Host, cfg.ServerGRPC.Port)))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := serverHTTP.Start(); err != nil {
			logg.Error("error starting HTTPServer",
				slog.String("address http", net.JoinHostPort(cfg.ServerHTTP.Host, cfg.ServerHTTP.Port)))
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if err := serverGRPC.Start(cfg.ServerGRPC); err != nil {
			logg.Error("error starting GRPCServer",
				slog.String("address grpc", net.JoinHostPort(cfg.ServerGRPC.Host, cfg.ServerGRPC.Port)))
			cancel()
		}
	}()

	wg.Wait()
}
