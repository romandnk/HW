package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv" //nolint:depguard
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/http"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/postgres"
	"github.com/spf13/viper"
)

var (
	ErrLoggerLevel                            = errors.New("invalid logger level")
	ErrLoggerRepresentation                   = errors.New("invalid logger representation")
	ErrLoggerEmptyLogFilePath                 = errors.New("logger file path cannot be empty")
	ErrLoggerFileNotExist                     = errors.New("logger file does not exist")
	ErrHTTPServerHost                         = errors.New("serverHTTP host must not be empty")
	ErrHTTPServerPortNotNumber                = errors.New("serverHTTP port must be a number")
	ErrHTTPServerPortWrongNumber              = errors.New("serverHTTP port must be in the interval from 0 to 65535")
	ErrParseHTTPServerReadTimeout             = errors.New("invalid serverHTTP read timeout")
	ErrParseHTTPServerWriteTimeout            = errors.New("invalid serverHTTP write timeout")
	ErrDBHost                                 = errors.New("database host must not be empty")
	ErrDBPortNotNumber                        = errors.New("database port must be a number")
	ErrDBPortWrongNumber                      = errors.New("database port must be in the interval from 0 to 65535")
	ErrDBInvalidDBName                        = errors.New("database name must not be empty")
	ErrDBInvalidSSLMode                       = errors.New("invalid database ssl mode")
	ErrDBMaxConns                             = errors.New("database max conns must be greater than 0")
	ErrDBMinConns                             = errors.New("database min conns must be greater than 0")
	ErrDBIncompatibleMaxAndMinConns           = errors.New("database max conns must be greater or equal to min conns")
	ErrParseMaxConnLifetime                   = errors.New("database errors parse MaxConnLifetime")
	ErrParseMaxConnIdleTime                   = errors.New("database errors parse MaxConnIdleTime")
	ErrHTTPServerReadTimeoutNotPositive       = errors.New("serverHTTP read timeout must be greater than 0")
	ErrHTTPServerWriteTimeoutNotPositive      = errors.New("serverHTTP write timeout must be greater than 0")
	ErrDBMaxConnLifetimeNotPositive           = errors.New("database MaxConnLifetime must be greater than 0")
	ErrDBMaxConnIdleTimeNotPositive           = errors.New("database MaxConnIdleTime must be greater than 0")
	ErrParseGRPCServerMaxConnectionIdle       = errors.New("invalid serverGRPC read max connection idle")
	ErrParseGRPCServerMaxConnectionAge        = errors.New("invalid serverGRPC read max connection age")
	ErrParseGRPCServerTime                    = errors.New("invalid serverGRPC read time")
	ErrParseGRPCServerTimeout                 = errors.New("invalid serverGRPC read timeout")
	ErrGRPCServerMaxConnectionAgeNotPositive  = errors.New("serverGRPC max connection age cannot be negative")
	ErrGRPCServerMaxConnectionIdleNotPositive = errors.New("serverGRPC max connection idle cannot be negative")
	ErrGRPCServerTimeoutNotPositive           = errors.New("serverGRPC timeout cannot be negative")
	ErrGRPCServerTimeNotPositive              = errors.New("serverGRPC time cannot be negative")
	ErrGRPCServerHost                         = errors.New("serverGRPC host must not be empty")
	ErrGRPCServerPortWrongNumber              = errors.New("serverGRPC port must be a number")
)

type Config struct {
	Logger      logger.Config
	ServerHTTP  internalhttp.Config
	ServerGRPC  grpc.Config
	StorageType string
	Storage     postgres.Config
}

func NewConfig(path string) (*Config, error) {
	viper.SetConfigFile(path) // find config file with specific path

	err := viper.ReadInConfig() // read config file
	if err != nil {
		return nil, fmt.Errorf("errors reading calendar config file: %w", err)
	}

	if err := godotenv.Load("./configs/calendar.env"); err != nil { // load .env into system
		return nil, fmt.Errorf("errors loading calendar.env: %w", err)
	}

	viper.SetEnvPrefix("calendar") // out env variables will look like CALENDAR_PASSWORD=password
	viper.AutomaticEnv()           // read env variables

	log := newLoggerConfig()
	err = validateLoggerConfig(log)
	if err != nil {
		return nil, err
	}

	serverHTTP, err := newServerHTTPConfig()
	if err != nil {
		return nil, err
	}
	err = validateServerHTTPConfig(serverHTTP)
	if err != nil {
		return nil, err
	}

	serverGRPC, err := newServerGRPCConfig()
	if err != nil {
		return nil, err
	}
	err = validateServerGRPCConfig(serverGRPC)
	if err != nil {
		return nil, err
	}

	storage, err := newStoragePostgresConfig()
	if err != nil {
		return nil, err
	}
	err = validateStoragePostgresConfig(storage)
	if err != nil {
		return nil, err
	}

	storageType := viper.GetString("storage.type")

	config := Config{
		Logger:      log,
		ServerHTTP:  serverHTTP,
		ServerGRPC:  serverGRPC,
		StorageType: storageType,
		Storage:     storage,
	}

	return &config, nil
}

func newLoggerConfig() logger.Config {
	level := viper.GetString("logger.level")
	representation := viper.GetString("logger.representation")
	lofFilePath := viper.GetString("logger.logs_file_path")
	return logger.Config{
		Level:          level,
		Representation: representation,
		LogFilePath:    lofFilePath,
	}
}

func newServerHTTPConfig() (internalhttp.Config, error) {
	host := viper.GetString("server_http.host")
	port := viper.GetString("server_http.port")

	readTimeoutStr := viper.GetString("server_http.read_timeout")
	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil {
		return internalhttp.Config{}, ErrParseHTTPServerReadTimeout
	}

	writeTimeoutStr := viper.GetString("server_http.write_timeout")
	writeTimeout, err := time.ParseDuration(writeTimeoutStr)
	if err != nil {
		return internalhttp.Config{}, ErrParseHTTPServerWriteTimeout
	}

	return internalhttp.Config{
		Host:         host,
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}, nil
}

func newServerGRPCConfig() (grpc.Config, error) {
	host := viper.GetString("server_grpc.host")
	port := viper.GetString("server_grpc.port")

	maxConnectionIdleStr := viper.GetString("server_grpc.max_connection_idle")
	maxConnectionIdle, err := time.ParseDuration(maxConnectionIdleStr)
	if err != nil {
		return grpc.Config{}, ErrParseGRPCServerMaxConnectionIdle
	}

	maxConnectionAgeStr := viper.GetString("server_grpc.max_connection_age")
	maxConnectionAge, err := time.ParseDuration(maxConnectionAgeStr)
	if err != nil {
		return grpc.Config{}, ErrParseGRPCServerMaxConnectionAge
	}

	timeStr := viper.GetString("server_grpc.time")
	parsedTime, err := time.ParseDuration(timeStr)
	if err != nil {
		return grpc.Config{}, ErrParseGRPCServerTime
	}

	timeoutStr := viper.GetString("server_grpc.timeout")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return grpc.Config{}, ErrParseGRPCServerTimeout
	}

	return grpc.Config{
		Host:              host,
		Port:              port,
		MaxConnectionIdle: maxConnectionIdle,
		MaxConnectionAge:  maxConnectionAge,
		Time:              parsedTime,
		Timeout:           timeout,
	}, nil
}

func newStoragePostgresConfig() (postgres.Config, error) {
	host := viper.GetString("storage.postgres.host")
	port := viper.GetString("storage.postgres.port")
	username := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("storage.postgres.db_name")
	sslmode := viper.GetString("storage.postgres.sslmode")
	maxConns := viper.GetInt("storage.postgres.max_conns")
	minConns := viper.GetInt("storage.postgres.min_conns")

	maxConnLifetimeStr := viper.GetString("storage.postgres.max_conn_lifetime")
	maxConnLifetime, err := time.ParseDuration(maxConnLifetimeStr)
	if err != nil {
		return postgres.Config{}, ErrParseMaxConnLifetime
	}
	maxConnIdleTimeStr := viper.GetString("storage.postgres.max_conn_idle_time")
	maxConnIdleTime, err := time.ParseDuration(maxConnIdleTimeStr)
	if err != nil {
		return postgres.Config{}, ErrParseMaxConnIdleTime
	}
	return postgres.Config{
		Host:            host,
		Port:            port,
		Username:        username,
		Password:        password,
		DBName:          dbName,
		SSLMode:         sslmode,
		MaxConns:        maxConns,
		MinConns:        minConns,
		MaxConnLifetime: maxConnLifetime,
		MaxConnIdleTime: maxConnIdleTime,
	}, nil
}

func validateLoggerConfig(l logger.Config) error {
	loggerLevels := map[string]struct{}{"INFO": {}, "DEBUG": {}, "ERROR": {}, "WARN": {}}
	if _, ok := loggerLevels[l.Level]; !ok {
		return ErrLoggerLevel
	}
	loggerRepresentations := map[string]struct{}{"JSON": {}, "TEXT": {}}
	if _, ok := loggerRepresentations[l.Representation]; !ok {
		return ErrLoggerRepresentation
	}
	if l.LogFilePath == "" {
		return ErrLoggerEmptyLogFilePath
	}
	_, err := os.Stat(l.LogFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrLoggerFileNotExist
		}
		return err
	}

	return nil
}

func validateServerHTTPConfig(s internalhttp.Config) error {
	if s.Host == "" {
		return ErrHTTPServerHost
	}
	port, err := strconv.Atoi(s.Port)
	if err != nil {
		return ErrHTTPServerPortNotNumber
	}
	if port < 0 || port > 65535 {
		return ErrHTTPServerPortWrongNumber
	}
	if s.ReadTimeout <= 0 {
		return ErrHTTPServerReadTimeoutNotPositive
	}
	if s.ReadTimeout <= 0 {
		return ErrHTTPServerWriteTimeoutNotPositive
	}

	return nil
}

func validateServerGRPCConfig(s grpc.Config) error {
	if s.Host == "" {
		return ErrGRPCServerHost
	}
	port, err := strconv.Atoi(s.Port)
	if err != nil {
		return err
	}
	if port < 0 || port > 65535 {
		return ErrGRPCServerPortWrongNumber
	}
	if s.MaxConnectionAge < 0 {
		return ErrGRPCServerMaxConnectionAgeNotPositive
	}
	if s.MaxConnectionIdle < 0 {
		return ErrGRPCServerMaxConnectionIdleNotPositive
	}
	if s.Timeout < 0 {
		return ErrGRPCServerTimeoutNotPositive
	}
	if s.Time < 0 {
		return ErrGRPCServerTimeNotPositive
	}

	return nil
}

func validateStoragePostgresConfig(st postgres.Config) error {
	if st.Host == "" {
		return ErrDBHost
	}
	port, err := strconv.Atoi(st.Port)
	if err != nil {
		return ErrDBPortNotNumber
	}
	if port < 0 || port > 65535 {
		return ErrDBPortWrongNumber
	}
	if st.DBName == "" {
		return ErrDBInvalidDBName
	}
	sslTypes := map[string]struct{}{"disable": {}, "verify-ca": {}, "require": {}, "verify-full": {}}
	if _, ok := sslTypes[st.SSLMode]; !ok {
		return ErrDBInvalidSSLMode
	}
	if st.MaxConns <= 0 {
		return ErrDBMaxConns
	}
	if st.MinConns <= 0 {
		return ErrDBMinConns
	}
	if st.MaxConns < st.MinConns {
		return ErrDBIncompatibleMaxAndMinConns
	}
	if st.MaxConnLifetime <= 0 {
		return ErrDBMaxConnLifetimeNotPositive
	}
	if st.MaxConnIdleTime <= 0 {
		return ErrDBMaxConnIdleTimeNotPositive
	}

	return nil
}
