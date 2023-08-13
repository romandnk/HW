package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"time"
)

var (
	ErrRabbitSenderEmptyUsername     = errors.New("username cannot be empty")
	ErrRabbitSenderEmptyPassword     = errors.New("password cannot be empty")
	ErrRabbitSenderEmptyHost         = errors.New("host cannot be empty")
	ErrRabbitSenderInvalidPort       = errors.New("port must be between 0 and 65535")
	ErrRabbitSenderNegativeHeartbeat = errors.New("heartbeat cannot be negative")
	ErrRabbitSenderEmptyExchangeName = errors.New("exchangeName cannot be empty")
	ErrRabbitSenderEmptyExchangeType = errors.New("exchangeType cannot be empty")
	ErrRabbitSenderEmptyQueueName    = errors.New("queueName cannot be empty")
	ErrRabbitSenderEmptyRoutingKey   = errors.New("routingKey cannot be empty")
	ErrRabbitSenderEmptyTag          = errors.New("tag cannot be empty")
)

type SenderRabbit struct {
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

func NewSenderRabbitConfig(path string) (SenderRabbit, error) {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return SenderRabbit{}, fmt.Errorf("errors reading rabbit config file: %w", err)
	}

	if err := godotenv.Load("./configs/sender.env"); err != nil {
		return SenderRabbit{}, fmt.Errorf("errors loading scheduler.env: %w", err)
	}

	viper.SetEnvPrefix("sender_rabbit")
	viper.AutomaticEnv()

	username := viper.GetString("user")
	password := viper.GetString("password")
	host := viper.GetString("rabbit_sender.host")
	port := viper.GetInt("rabbit_sender.port")
	heartbeat := viper.GetDuration("rabbit_sender.heartbeat")
	exchangeName := viper.GetString("rabbit_sender.exchange_name")
	exchangeType := viper.GetString("rabbit_sender.exchange_type")
	durableExchange := viper.GetBool("rabbit_sender.durable_exchange")
	autoDeleteExchange := viper.GetBool("rabbit_sender.autoDelete_exchange")
	queueName := viper.GetString("rabbit_sender.queue_name")
	durableQueue := viper.GetBool("rabbit_sender.durable_queue")
	autoDeleteQueue := viper.GetBool("rabbit_sender.autoDelete_queue")
	routingKey := viper.GetString("rabbit_sender.touting_key")
	tag := viper.GetString("rabbit_sender.tag")

	config := SenderRabbit{
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
	}

	err = validateSenderRabbit(config)
	if err != nil {
		return SenderRabbit{}, err
	}

	return config, nil
}

func validateSenderRabbit(s SenderRabbit) error {
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
