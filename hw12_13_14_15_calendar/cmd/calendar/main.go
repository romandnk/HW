package main

import (
	"context"
	"flag"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/cmd/config"
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

	cfg, err := config.NewCalendarConfig(configFile)
	if err != nil {
		log.Fatalf("calendar config error: %s", err.Error())
	}

	logg := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Representation)

	logg.Info("use logging")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var st storage.Storage

	// use memory storage or sql storage
	switch cfg.Storage.Type {
	case "memory":
		st = memorystorage.NewStorageMemory()
	case "postgres":
		db, err := sqlstorage.NewPostgresDB(ctx, cfg.Storage.DB)
		if err != nil {
			logg.Error("errors connecting db",
				slog.String("errors", err.Error()),
				slog.String("address", cfg.Storage.DB.Host+":"+cfg.Storage.DB.Port))
			cancel()
		}
		defer db.Close()

		st = sqlstorage.NewStorageSQL(db)
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
			logg.Error("errors stopping HTTPServer",
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
			logg.Error("errors starting HTTPServer",
				slog.String("address http", net.JoinHostPort(cfg.ServerHTTP.Host, cfg.ServerHTTP.Port)))
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if err := serverGRPC.Start(cfg.ServerGRPC); err != nil {
			logg.Error("errors starting GRPCServer",
				slog.String("address grpc", net.JoinHostPort(cfg.ServerGRPC.Host, cfg.ServerGRPC.Port)))
			cancel()
		}
	}()

	wg.Wait()
}
