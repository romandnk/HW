package rmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"time"
)

type RabbitConfig struct {
	Username  string
	Password  string
	Host      string
	Port      int
	Heartbeat time.Duration
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
	log     logger.Logger
}

func NewConsumer(log logger.Logger, cfg RabbitConfig) (*Consumer, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)

	config := amqp.Config{
		Heartbeat: cfg.Heartbeat,
	}

	connection, err := amqp.DialConfig(url, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{conn: connection, log: log}, nil
}
