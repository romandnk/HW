package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv" //nolint:depguard
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/http"
	dbconf "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/viper"
)

var (
	ErrLoggerLevel                            = errors.New("invalid logger level")
	ErrLoggerRepresentation                   = errors.New("invalid logger representation")
	ErrLoggerEmptyLogFilePath                 = errors.New("logger file path cannot be empty")
	ErrLoggerFileNotExist                     = errors.New("logger file does not existgit ")
	ErrHTTPServerHost                         = errors.New("serverHTTP host must not be empty")
	ErrHTTPServerPortNotNumber                = errors.New("serverHTTP port must be a number")
	ErrHTTPServerPortWrongNumber              = errors.New("serverHTTP port must be in the interval from 0 to 65535")
	ErrParseHTTPServerReadTimeout             = errors.New("invalid serverHTTP read timeout")
	ErrParseHTTPServerWriteTimeout            = errors.New("invalid serverHTTP write timeout")
	ErrInvalidStorageType                     = errors.New("invalid storage type")
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
)

var (
	memStorage      = "memory"
	postgresStorage = "postgres"
)

type CalendarConfig struct {
	Logger     LoggerConf
	ServerHTTP internalhttp.ServerHTTPConfig
	ServerGRPC grpc.ServerGRPCConfig
	Storage    StorageConf
}

type LoggerConf struct {
	Level          string
	Representation string
	LogFilePath    string
}

type StorageConf struct {
	Type string
	DB   dbconf.DBConf
}

func NewCalendarConfig(path string) (*CalendarConfig, error) {
	viper.SetConfigFile(path) // find config file with specific path

	err := viper.ReadInConfig() // read config file
	if err != nil {
		return &CalendarConfig{}, fmt.Errorf("errors reading config file: %w", err)
	}

	if err := godotenv.Load("./configs/.env"); err != nil { // load .env into system
		return &CalendarConfig{}, fmt.Errorf("errors loading .env: %w", err)
	}

	viper.SetEnvPrefix("calendar") // out env variables will look like CALENDAR_PASSWORD=password
	viper.AutomaticEnv()           // read env variables

	logger := newLoggerConf()
	err = validateLogger(logger)
	if err != nil {
		return &CalendarConfig{}, err
	}

	serverHTTP, err := newServerHTTPConf()
	if err != nil {
		return &CalendarConfig{}, err
	}
	err = validateServerHTTP(serverHTTP)
	if err != nil {
		return &CalendarConfig{}, err
	}

	serverGRPC, err := newServerGRPCConf()
	if err != nil {
		return &CalendarConfig{}, err
	}
	err = validateServerGRPC(serverGRPC)
	if err != nil {
		return &CalendarConfig{}, err
	}

	storage, err := newStorageConf()
	if err != nil {
		return &CalendarConfig{}, err
	}
	err = validateStorage(storage)
	if err != nil {
		return &CalendarConfig{}, err
	}

	config := CalendarConfig{
		Logger:     logger,
		ServerHTTP: serverHTTP,
		ServerGRPC: serverGRPC,
		Storage:    storage,
	}

	return &config, nil
}

func newLoggerConf() LoggerConf {
	level := viper.GetString("logger.level")
	representation := viper.GetString("logger.representation")
	lofFilePath := viper.GetString("logger.logs_file_path")
	return LoggerConf{
		Level:          level,
		Representation: representation,
		LogFilePath:    lofFilePath,
	}
}

func newServerHTTPConf() (internalhttp.ServerHTTPConfig, error) {
	host := viper.GetString("server_http.host")
	port := viper.GetString("server_http.port")

	readTimeoutStr := viper.GetString("server_http.read_timeout")
	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil {
		return internalhttp.ServerHTTPConfig{}, ErrParseHTTPServerReadTimeout
	}

	writeTimeoutStr := viper.GetString("server_http.write_timeout")
	writeTimeout, err := time.ParseDuration(writeTimeoutStr)
	if err != nil {
		return internalhttp.ServerHTTPConfig{}, ErrParseHTTPServerWriteTimeout
	}

	return internalhttp.ServerHTTPConfig{
		Host:         host,
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}, nil
}

func newServerGRPCConf() (grpc.ServerGRPCConfig, error) {
	host := viper.GetString("server_grpc.host")
	port := viper.GetString("server_grpc.port")

	maxConnectionIdleStr := viper.GetString("server_grpc.max_connection_idle")
	maxConnectionIdle, err := time.ParseDuration(maxConnectionIdleStr)
	if err != nil {
		return grpc.ServerGRPCConfig{}, ErrParseGRPCServerMaxConnectionIdle
	}

	maxConnectionAgeStr := viper.GetString("server_grpc.max_connection_age")
	maxConnectionAge, err := time.ParseDuration(maxConnectionAgeStr)
	if err != nil {
		return grpc.ServerGRPCConfig{}, ErrParseGRPCServerMaxConnectionAge
	}

	timeStr := viper.GetString("server_grpc.time")
	parsedTime, err := time.ParseDuration(timeStr)
	if err != nil {
		return grpc.ServerGRPCConfig{}, ErrParseGRPCServerTime
	}

	timeoutStr := viper.GetString("server_grpc.timeout")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return grpc.ServerGRPCConfig{}, ErrParseGRPCServerTimeout
	}

	return grpc.ServerGRPCConfig{
		Host:              host,
		Port:              port,
		MaxConnectionIdle: maxConnectionIdle,
		MaxConnectionAge:  maxConnectionAge,
		Time:              parsedTime,
		Timeout:           timeout,
	}, nil
}

func newStorageConf() (StorageConf, error) {
	storageType := viper.GetString("storage.type")
	switch storageType {
	case memStorage:
		return StorageConf{
			Type: memStorage,
		}, nil
	case postgresStorage:
		host := viper.GetString("storage.database.host")
		port := viper.GetString("storage.database.port")
		username := viper.GetString("DB_USER")
		password := viper.GetString("DB_PASSWORD")
		dbName := viper.GetString("storage.database.db_name")
		sslmode := viper.GetString("storage.database.sslmode")
		maxConns := viper.GetInt("storage.database.max_conns")
		minConns := viper.GetInt("storage.database.min_conns")

		maxConnLifetimeStr := viper.GetString("storage.database.max_conn_lifetime")
		maxConnLifetime, err := time.ParseDuration(maxConnLifetimeStr)
		if err != nil {
			return StorageConf{}, ErrParseMaxConnLifetime
		}

		maxConnIdleTimeStr := viper.GetString("storage.database.max_conn_idle_time")
		maxConnIdleTime, err := time.ParseDuration(maxConnIdleTimeStr)
		if err != nil {
			return StorageConf{}, ErrParseMaxConnIdleTime
		}

		DBconf := dbconf.DBConf{
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
		}
		return StorageConf{
			Type: postgresStorage,
			DB:   DBconf,
		}, nil
	default:
		return StorageConf{}, ErrInvalidStorageType
	}
}

func validateLogger(l LoggerConf) error {
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

func validateServerHTTP(s internalhttp.ServerHTTPConfig) error {
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

func validateServerGRPC(s grpc.ServerGRPCConfig) error {
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

func validateStorage(st StorageConf) error {
	if st.Type == postgresStorage { //nolint:nestif
		if st.DB.Host == "" {
			return ErrDBHost
		}
		port, err := strconv.Atoi(st.DB.Port)
		if err != nil {
			return ErrDBPortNotNumber
		}
		if port < 0 || port > 65535 {
			return ErrDBPortWrongNumber
		}

		if st.DB.DBName == "" {
			return ErrDBInvalidDBName
		}
		sslTypes := map[string]struct{}{"disable": {}, "verify-ca": {}, "require": {}, "verify-full": {}}
		if _, ok := sslTypes[st.DB.SSLMode]; !ok {
			return ErrDBInvalidSSLMode
		}
		if st.DB.MaxConns <= 0 {
			return ErrDBMaxConns
		}
		if st.DB.MinConns <= 0 {
			return ErrDBMinConns
		}
		if st.DB.MaxConns < st.DB.MinConns {
			return ErrDBIncompatibleMaxAndMinConns
		}
		if st.DB.MaxConnLifetime <= 0 {
			return ErrDBMaxConnLifetimeNotPositive
		}
		if st.DB.MaxConnIdleTime <= 0 {
			return ErrDBMaxConnIdleTimeNotPositive
		}
	}

	return nil
}
