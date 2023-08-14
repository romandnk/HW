package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/joho/godotenv"
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
)

type SchedulerConfig struct {
	MQ     RabbitSchedulerConfig
	Logger LoggerSchedulerConfig
}

type RabbitSchedulerConfig struct {
	Username           string
	Password           string
	Host               string
	Port               int
	Heartbeat          time.Duration
	ExchangeName       string
	QueueName          string
	ExchangeType       string
	DurableExchange    bool
	DurableQueue       bool
	AutoDeleteExchange bool
	AutoDeleteQueue    bool
	RoutingKey         string
	DeliveryMode       int
}

type LoggerSchedulerConfig struct {
	Level          string
	Representation string
}

func NewSchedulerConfig(path string) (*SchedulerConfig, error) {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("errors reading rabbit scheduler config file: %w", err)
	}

	if err := godotenv.Load("./configs/scheduler.env"); err != nil {
		return nil, fmt.Errorf("errors loading scheduler.env: %w", err)
	}

	viper.SetEnvPrefix("scheduler_rabbit")
	viper.AutomaticEnv()

	rabbitConfig, err := newRabbitSchedulerConfig()
	if err != nil {
		return nil, err
	}
	err = validateRabbitSchedulerConfig(rabbitConfig)
	if err != nil {
		return nil, err
	}

	logger := newLoggerSchedulerConfig()
	err = validateLoggerSchedulerConfig(logger)
	if err != nil {
		return nil, err
	}

	config := SchedulerConfig{
		MQ:     rabbitConfig,
		Logger: logger,
	}

	return &config, nil
}

func newRabbitSchedulerConfig() (RabbitSchedulerConfig, error) {
	username := viper.GetString("user")
	password := viper.GetString("password")
	host := viper.GetString("rabbit_scheduler.host")
	port := viper.GetInt("rabbit_scheduler.port")

	heartbeatStr := viper.GetString("rabbit_scheduler.heartbeat")
	heartbeat, err := time.ParseDuration(heartbeatStr)
	if err != nil {
		return RabbitSchedulerConfig{}, ErrRabbitSchedulerParseHeartbeat
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

	return RabbitSchedulerConfig{
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

func validateRabbitSchedulerConfig(s RabbitSchedulerConfig) error {
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

func newLoggerSchedulerConfig() LoggerSchedulerConfig {
	level := viper.GetString("logger.level")
	representation := viper.GetString("logger.representation")
	return LoggerSchedulerConfig{
		Level:          level,
		Representation: representation,
	}
}

func validateLoggerSchedulerConfig(l LoggerSchedulerConfig) error {
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
