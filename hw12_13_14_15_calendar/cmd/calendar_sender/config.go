package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/mq/rabbitmq"
	"github.com/spf13/viper"
)

var (
	ErrRabbitSenderEmptyUsername        = errors.New("username cannot be empty")
	ErrRabbitSenderEmptyPassword        = errors.New("password cannot be empty")
	ErrRabbitSenderEmptyHost            = errors.New("host cannot be empty")
	ErrRabbitSenderInvalidPort          = errors.New("port must be between 0 and 65535")
	ErrRabbitSenderNegativeHeartbeat    = errors.New("heartbeat cannot be negative")
	ErrRabbitSenderEmptyExchangeName    = errors.New("exchangeName cannot be empty")
	ErrRabbitSenderEmptyExchangeType    = errors.New("exchangeType cannot be empty")
	ErrRabbitSenderEmptyQueueName       = errors.New("queueName cannot be empty")
	ErrRabbitSenderEmptyRoutingKey      = errors.New("routingKey cannot be empty")
	ErrRabbitSenderEmptyTag             = errors.New("tag cannot be empty")
	ErrRabbitSenderConfigParseHeartbeat = errors.New("invalid heartbeat")
	ErrRabbitSenderInvalidExchangeType  = errors.New("invalid exchange type: direct, fanout, headers, topic")
	ErrLoggerLevel                      = errors.New("invalid logger level")
	ErrLoggerRepresentation             = errors.New("invalid logger representation")
)

type Config struct {
	MQ     rabbitmq.SenderConfig
	Logger logger.Config
}

func NewConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("errors reading rabbit sender config file: %w", err)
	}

	if err := godotenv.Load("./configs/sender.env"); err != nil {
		return nil, fmt.Errorf("errors loading sender.env: %w", err)
	}

	viper.SetEnvPrefix("sender_rabbit")
	viper.AutomaticEnv()

	rabbitConfig, err := newRabbitSenderConfig()
	if err != nil {
		return nil, err
	}
	err = validateRabbitSenderConfig(rabbitConfig)
	if err != nil {
		return nil, err
	}

	log := newLoggerConfig()
	err = validateLoggerConfig(log)
	if err != nil {
		return nil, err
	}

	config := Config{
		MQ:     rabbitConfig,
		Logger: log,
	}

	return &config, nil
}

func newRabbitSenderConfig() (rabbitmq.SenderConfig, error) {
	username := viper.GetString("user")
	password := viper.GetString("password")
	host := viper.GetString("rabbit_sender.host")
	port := viper.GetInt("rabbit_sender.port")

	heartbeatStr := viper.GetString("rabbit_sender.heartbeat")
	heartbeat, err := time.ParseDuration(heartbeatStr)
	if err != nil {
		return rabbitmq.SenderConfig{}, ErrRabbitSenderConfigParseHeartbeat
	}

	exchangeName := viper.GetString("rabbit_sender.exchange_name")
	exchangeType := viper.GetString("rabbit_sender.exchange_type")
	durableExchange := viper.GetBool("rabbit_sender.durable_exchange")
	autoDeleteExchange := viper.GetBool("rabbit_sender.auto_delete_exchange")
	queueName := viper.GetString("rabbit_sender.queue_name")
	durableQueue := viper.GetBool("rabbit_sender.durable_queue")
	autoDeleteQueue := viper.GetBool("rabbit_sender.auto_delete_queue")
	routingKey := viper.GetString("rabbit_sender.routing_key")
	tag := viper.GetString("rabbit_sender.tag")

	return rabbitmq.SenderConfig{
		Username:           username,
		Password:           password,
		Host:               host,
		Port:               port,
		Heartbeat:          heartbeat,
		ExchangeName:       exchangeName,
		ExchangeType:       exchangeType,
		DurableExchange:    durableExchange,
		AutoDeleteExchange: autoDeleteExchange,
		QueueName:          queueName,
		DurableQueue:       durableQueue,
		AutoDeleteQueue:    autoDeleteQueue,
		RoutingKey:         routingKey,
		Tag:                tag,
	}, nil
}

func validateRabbitSenderConfig(s rabbitmq.SenderConfig) error {
	if s.Username == "" {
		return ErrRabbitSenderEmptyUsername
	}
	if s.Password == "" {
		return ErrRabbitSenderEmptyPassword
	}
	if s.Host == "" {
		return ErrRabbitSenderEmptyHost
	}
	if s.Port < 0 || s.Port > 65535 {
		return ErrRabbitSenderInvalidPort
	}
	if s.Heartbeat < 0 {
		return ErrRabbitSenderNegativeHeartbeat
	}
	if s.ExchangeName == "" {
		return ErrRabbitSenderEmptyExchangeName
	}
	if s.ExchangeType == "" {
		return ErrRabbitSenderEmptyExchangeType
	}
	exchangeTypes := map[string]struct{}{"direct": {}, "fanout": {}, "topic": {}, "headers": {}}
	if _, ok := exchangeTypes[s.ExchangeType]; !ok {
		return ErrRabbitSenderInvalidExchangeType
	}
	if s.QueueName == "" {
		return ErrRabbitSenderEmptyQueueName
	}
	if s.RoutingKey == "" {
		return ErrRabbitSenderEmptyRoutingKey
	}
	if s.Tag == "" {
		return ErrRabbitSenderEmptyTag
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
