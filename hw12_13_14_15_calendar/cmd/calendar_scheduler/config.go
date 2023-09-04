package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/mq/rabbitmq"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/postgres"
	"github.com/spf13/viper"
)

var (
	ErrRabbitSchedulerParseHeartbeat      = errors.New("invalid heartbeat")
	ErrRabbitSchedulerEmptyUsername       = errors.New("username cannot be empty")
	ErrRabbitSchedulerEmptyPassword       = errors.New("password cannot be empty")
	ErrRabbitSchedulerEmptyHost           = errors.New("host cannot be empty")
	ErrRabbitSchedulerInvalidPort         = errors.New("port must be between 0 and 65535")
	ErrRabbitSchedulerNegativeHeartbeat   = errors.New("heartbeat cannot be negative")
	ErrRabbitSchedulerEmptyExchangeName   = errors.New("exchangeName cannot be empty")
	ErrRabbitSchedulerEmptyExchangeType   = errors.New("exchangeType cannot be empty")
	ErrRabbitSchedulerInvalidExchangeType = errors.New("invalid exchange type: direct, fanout, headers, topic")
	ErrRabbitSchedulerEmptyRoutingKey     = errors.New("routingKey cannot be empty")
	ErrRabbitSchedulerInvalidDeliveryMode = errors.New("invalid delivery mode: Transient (0 or 1) or Persistent (2)")
	ErrRabbitSchedulerEmptyQueueName      = errors.New("queueName cannot be empty")
	ErrLoggerLevel                        = errors.New("invalid logger level")
	ErrLoggerRepresentation               = errors.New("invalid logger representation")
	ErrDBHost                             = errors.New("database host must not be empty")
	ErrDBPortNotNumber                    = errors.New("database port must be a number")
	ErrDBPortWrongNumber                  = errors.New("database port must be in the interval from 0 to 65535")
	ErrDBInvalidDBName                    = errors.New("database name must not be empty")
	ErrDBInvalidSSLMode                   = errors.New("invalid database ssl mode")
	ErrDBMaxConns                         = errors.New("database max conns must be greater than 0")
	ErrDBMinConns                         = errors.New("database min conns must be greater than 0")
	ErrDBIncompatibleMaxAndMinConns       = errors.New("database max conns must be greater or equal to min conns")
	ErrParseMaxConnLifetime               = errors.New("database errors parse MaxConnLifetime")
	ErrParseMaxConnIdleTime               = errors.New("database errors parse MaxConnIdleTime")
	ErrDBMaxConnLifetimeNotPositive       = errors.New("database MaxConnLifetime must be greater than 0")
	ErrDBMaxConnIdleTimeNotPositive       = errors.New("database MaxConnIdleTime must be greater than 0")
)

type Config struct {
	MQ                   rabbitmq.ProducerConfig
	Logger               logger.Config
	StorageType          string
	Storage              postgres.Config
	TimeToSchedule       time.Duration
	TimeToDeleteOutdated time.Duration
}

func NewConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("errors reading rabbit scheduler config file: %w", err)
	}

	viper.SetEnvPrefix("scheduler")
	viper.AutomaticEnv()

	rabbitConfig, err := newRabbitSchedulerConfig()
	if err != nil {
		return nil, err
	}
	err = validateRabbitSchedulerConfig(rabbitConfig)
	if err != nil {
		return nil, err
	}

	log := newLoggerConfig()
	err = validateLoggerConfig(log)
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

	timeToScheduleStr := viper.GetString("general_preferences.time_to_schedule")
	timeToSchedule, err := time.ParseDuration(timeToScheduleStr)
	if err != nil {
		return nil, err
	}

	timeToDeleteOutdatedStr := viper.GetString("general_preferences.time_to_delete_outdated")
	timeToDeleteOutdated, err := time.ParseDuration(timeToDeleteOutdatedStr)
	if err != nil {
		return nil, err
	}

	config := Config{
		MQ:                   rabbitConfig,
		Logger:               log,
		StorageType:          storageType,
		Storage:              storage,
		TimeToSchedule:       timeToSchedule,
		TimeToDeleteOutdated: timeToDeleteOutdated,
	}

	return &config, nil
}

func newRabbitSchedulerConfig() (rabbitmq.ProducerConfig, error) {
	username := viper.GetString("rabbit_user")
	password := viper.GetString("rabbit_password")
	host := viper.GetString("rabbit_scheduler.host")
	port := viper.GetInt("rabbit_scheduler.port")

	heartbeatStr := viper.GetString("rabbit_scheduler.heartbeat")
	heartbeat, err := time.ParseDuration(heartbeatStr)
	if err != nil {
		return rabbitmq.ProducerConfig{}, ErrRabbitSchedulerParseHeartbeat
	}

	exchangeName := viper.GetString("rabbit_scheduler.exchange_name")
	queueName := viper.GetString("rabbit_scheduler.queue_name")
	exchangeType := viper.GetString("rabbit_scheduler.exchange_type")
	durableExchange := viper.GetBool("rabbit_scheduler.durable_exchange")
	durableQueue := viper.GetBool("rabbit_scheduler.durable_queue")
	autoDeleteExchange := viper.GetBool("rabbit_scheduler.auto_delete_exchange")
	autoDeleteQueue := viper.GetBool("rabbit_scheduler.auto_delete_queue")
	routingKey := viper.GetString("rabbit_scheduler.routing_key")
	deliveryMode := viper.GetInt("rabbit_scheduler.delivery_mode")

	return rabbitmq.ProducerConfig{
		Username:           username,
		Password:           password,
		Host:               host,
		Port:               port,
		Heartbeat:          heartbeat,
		ExchangeName:       exchangeName,
		QueueName:          queueName,
		ExchangeType:       exchangeType,
		DurableExchange:    durableExchange,
		DurableQueue:       durableQueue,
		AutoDeleteExchange: autoDeleteExchange,
		AutoDeleteQueue:    autoDeleteQueue,
		RoutingKey:         routingKey,
		DeliveryMode:       deliveryMode,
	}, nil
}

func validateRabbitSchedulerConfig(s rabbitmq.ProducerConfig) error {
	if s.Username == "" {
		return ErrRabbitSchedulerEmptyUsername
	}
	if s.Password == "" {
		return ErrRabbitSchedulerEmptyPassword
	}
	if s.Host == "" {
		return ErrRabbitSchedulerEmptyHost
	}
	if s.Port < 0 || s.Port > 65535 {
		return ErrRabbitSchedulerInvalidPort
	}
	if s.Heartbeat < 0 {
		return ErrRabbitSchedulerNegativeHeartbeat
	}
	if s.ExchangeName == "" {
		return ErrRabbitSchedulerEmptyExchangeName
	}
	if s.QueueName == "" {
		return ErrRabbitSchedulerEmptyQueueName
	}
	if s.ExchangeType == "" {
		return ErrRabbitSchedulerEmptyExchangeType
	}
	exchangeTypes := map[string]struct{}{"direct": {}, "fanout": {}, "topic": {}, "headers": {}}
	if _, ok := exchangeTypes[s.ExchangeType]; !ok {
		return ErrRabbitSchedulerInvalidExchangeType
	}
	if s.RoutingKey == "" {
		return ErrRabbitSchedulerEmptyRoutingKey
	}
	if s.DeliveryMode < 0 || s.DeliveryMode > 2 {
		return ErrRabbitSchedulerInvalidDeliveryMode
	}
	return nil
}

func newLoggerConfig() logger.Config {
	level := viper.GetString("logger.level")
	representation := viper.GetString("logger.representation")
	return logger.Config{
		Level:          level,
		Representation: representation,
	}
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

	return nil
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
