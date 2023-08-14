package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	ErrRabbitSenderEmptyUsername       = errors.New("username cannot be empty")
	ErrRabbitSenderEmptyPassword       = errors.New("password cannot be empty")
	ErrRabbitSenderEmptyHost           = errors.New("host cannot be empty")
	ErrRabbitSenderInvalidPort         = errors.New("port must be between 0 and 65535")
	ErrRabbitSenderNegativeHeartbeat   = errors.New("heartbeat cannot be negative")
	ErrRabbitSenderEmptyExchangeName   = errors.New("exchangeName cannot be empty")
	ErrRabbitSenderEmptyExchangeType   = errors.New("exchangeType cannot be empty")
	ErrRabbitSenderEmptyQueueName      = errors.New("queueName cannot be empty")
	ErrRabbitSenderEmptyRoutingKey     = errors.New("routingKey cannot be empty")
	ErrRabbitSenderEmptyTag            = errors.New("tag cannot be empty")
	ErrRabbitConfigParseHeartbeat      = errors.New("invalid heartbeat")
	ErrRabbitSenderInvalidExchangeType = errors.New("invalid exchange type: direct, fanout, headers, topic")
)

type SenderConfig struct {
	MQ     RabbitConfig
	Logger LoggerSenderConfig
}

type RabbitConfig struct {
	Username           string
	Password           string
	Host               string
	Port               int
	Heartbeat          time.Duration
	ExchangeName       string
	ExchangeType       string
	DurableExchange    bool
	AutoDeleteExchange bool
	QueueName          string
	DurableQueue       bool
	AutoDeleteQueue    bool
	RoutingKey         string
	Tag                string
}

type LoggerSenderConfig struct {
	Level          string
	Representation string
}

func NewSenderConfig(path string) (*SenderConfig, error) {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("errors reading rabbit config file: %w", err)
	}

	if err := godotenv.Load("./configs/sender.env"); err != nil {
		return nil, fmt.Errorf("errors loading scheduler.env: %w", err)
	}

	viper.SetEnvPrefix("sender_rabbit")
	viper.AutomaticEnv()

	rabbitConfig, err := newRabbitConfig()
	if err != nil {
		return nil, err
	}
	err = validateSenderRabbit(rabbitConfig)
	if err != nil {
		return nil, err
	}

	logger := newLoggerSenderConf()
	err = validateSenderLogger(logger)
	if err != nil {
		return nil, err
	}

	config := SenderConfig{
		MQ:     rabbitConfig,
		Logger: logger,
	}

	return &config, nil
}

func newRabbitConfig() (RabbitConfig, error) {
	username := viper.GetString("user")
	password := viper.GetString("password")
	host := viper.GetString("rabbit_sender.host")
	port := viper.GetInt("rabbit_sender.port")

	heartbeatStr := viper.GetString("rabbit_sender.heartbeat")
	heartbeat, err := time.ParseDuration(heartbeatStr)
	if err != nil {
		return RabbitConfig{}, ErrRabbitConfigParseHeartbeat
	}

	exchangeName := viper.GetString("rabbit_sender.exchange_name")
	exchangeType := viper.GetString("rabbit_sender.exchange_type")
	durableExchange := viper.GetBool("rabbit_sender.durable_exchange")
	autoDeleteExchange := viper.GetBool("rabbit_sender.autoDelete_exchange")
	queueName := viper.GetString("rabbit_sender.queue_name")
	durableQueue := viper.GetBool("rabbit_sender.durable_queue")
	autoDeleteQueue := viper.GetBool("rabbit_sender.autoDelete_queue")
	routingKey := viper.GetString("rabbit_sender.touting_key")
	tag := viper.GetString("rabbit_sender.tag")

	return RabbitConfig{
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

func validateSenderRabbit(s RabbitConfig) error {
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

func newLoggerSenderConf() LoggerSenderConfig {
	level := viper.GetString("logger.level")
	representation := viper.GetString("logger.representation")
	return LoggerSenderConfig{
		Level:          level,
		Representation: representation,
	}
}

func validateSenderLogger(l LoggerSenderConfig) error {
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
