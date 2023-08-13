package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv" //nolint:depguard
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
	Logger     LoggerCalendarConfig
	ServerHTTP ServerHTTPConfig
	ServerGRPC ServerGRPCConfig
	Storage    StorageConfig
}

type LoggerCalendarConfig struct {
	Level          string
	Representation string
	LogFilePath    string
}

type ServerHTTPConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type ServerGRPCConfig struct {
	Host              string
	Port              string
	MaxConnectionIdle time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
	Timeout           time.Duration
}

type StorageConfig struct {
	Type string
	DB   DBConfig
}

type DBConfig struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	SSLMode         string
	MaxConns        int           // max connections in the pool
	MinConns        int           // min connections in the pool which must be opened
	MaxConnLifetime time.Duration // time after which db conn will be removed from the pool if there was no active use.
	MaxConnIdleTime time.Duration // time after which an inactive connection in the pool will be closed and deleted.
}

func NewCalendarConfig(path string) (*CalendarConfig, error) {
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

	logger := newLoggerCalendarConf()
	err = validateCalendarLogger(logger)
	if err != nil {
		return nil, err
	}

	serverHTTP, err := newServerHTTPConf()
	if err != nil {
		return nil, err
	}
	err = validateServerHTTP(serverHTTP)
	if err != nil {
		return nil, err
	}

	serverGRPC, err := newServerGRPCConf()
	if err != nil {
		return nil, err
	}
	err = validateServerGRPC(serverGRPC)
	if err != nil {
		return nil, err
	}

	storage, err := newStorageConf()
	if err != nil {
		return nil, err
	}
	err = validateStorage(storage)
	if err != nil {
		return nil, err
	}

	config := CalendarConfig{
		Logger:     logger,
		ServerHTTP: serverHTTP,
		ServerGRPC: serverGRPC,
		Storage:    storage,
	}

	return &config, nil
}

func newLoggerCalendarConf() LoggerCalendarConfig {
	level := viper.GetString("logger.level")
	representation := viper.GetString("logger.representation")
	lofFilePath := viper.GetString("logger.logs_file_path")
	return LoggerCalendarConfig{
		Level:          level,
		Representation: representation,
		LogFilePath:    lofFilePath,
	}
}

func newServerHTTPConf() (ServerHTTPConfig, error) {
	host := viper.GetString("server_http.host")
	port := viper.GetString("server_http.port")

	readTimeoutStr := viper.GetString("server_http.read_timeout")
	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil {
		return ServerHTTPConfig{}, ErrParseHTTPServerReadTimeout
	}

	writeTimeoutStr := viper.GetString("server_http.write_timeout")
	writeTimeout, err := time.ParseDuration(writeTimeoutStr)
	if err != nil {
		return ServerHTTPConfig{}, ErrParseHTTPServerWriteTimeout
	}

	return ServerHTTPConfig{
		Host:         host,
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}, nil
}

func newServerGRPCConf() (ServerGRPCConfig, error) {
	host := viper.GetString("server_grpc.host")
	port := viper.GetString("server_grpc.port")

	maxConnectionIdleStr := viper.GetString("server_grpc.max_connection_idle")
	maxConnectionIdle, err := time.ParseDuration(maxConnectionIdleStr)
	if err != nil {
		return ServerGRPCConfig{}, ErrParseGRPCServerMaxConnectionIdle
	}

	maxConnectionAgeStr := viper.GetString("server_grpc.max_connection_age")
	maxConnectionAge, err := time.ParseDuration(maxConnectionAgeStr)
	if err != nil {
		return ServerGRPCConfig{}, ErrParseGRPCServerMaxConnectionAge
	}

	timeStr := viper.GetString("server_grpc.time")
	parsedTime, err := time.ParseDuration(timeStr)
	if err != nil {
		return ServerGRPCConfig{}, ErrParseGRPCServerTime
	}

	timeoutStr := viper.GetString("server_grpc.timeout")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return ServerGRPCConfig{}, ErrParseGRPCServerTimeout
	}

	return ServerGRPCConfig{
		Host:              host,
		Port:              port,
		MaxConnectionIdle: maxConnectionIdle,
		MaxConnectionAge:  maxConnectionAge,
		Time:              parsedTime,
		Timeout:           timeout,
	}, nil
}

func newStorageConf() (StorageConfig, error) {
	storageType := viper.GetString("storage.type")
	switch storageType {
	case memStorage:
		return StorageConfig{
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
			return StorageConfig{}, ErrParseMaxConnLifetime
		}

		maxConnIdleTimeStr := viper.GetString("storage.database.max_conn_idle_time")
		maxConnIdleTime, err := time.ParseDuration(maxConnIdleTimeStr)
		if err != nil {
			return StorageConfig{}, ErrParseMaxConnIdleTime
		}

		DBconf := DBConfig{
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
		return StorageConfig{
			Type: postgresStorage,
			DB:   DBconf,
		}, nil
	default:
		return StorageConfig{}, ErrInvalidStorageType
	}
}

func validateCalendarLogger(l LoggerCalendarConfig) error {
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

func validateServerHTTP(s ServerHTTPConfig) error {
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

func validateServerGRPC(s ServerGRPCConfig) error {
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

func validateStorage(st StorageConfig) error {
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
