package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
)

var ErrSchedulerRabbitNilChannel = errors.New("rabbit scheduler: channel is nil")

type ProducerConfig struct {
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

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	log     logger.Logger
	cfg     ProducerConfig
}

func NewProducer(cfg ProducerConfig, log logger.Logger) (*Producer, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conf := amqp.Config{
		Heartbeat: cfg.Heartbeat,
	}
	log.Info("connecting to rabbit scheduler...")
	conn, err := amqp.DialConfig(url, conf)
	if err != nil {
		return nil, err
	}

	log.Info("connected to rabbit scheduler")

	log.Info("opening scheduler channel...")
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	log.Info("declaring scheduler exchange...")
	err = ch.ExchangeDeclare(
		cfg.ExchangeName,
		cfg.ExchangeType,
		cfg.DurableExchange,
		cfg.AutoDeleteExchange,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		cfg.QueueName,
		cfg.DurableQueue,
		cfg.AutoDeleteQueue,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		queue.Name,
		cfg.RoutingKey,
		cfg.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Producer{
		conn:    conn,
		channel: ch,
		log:     log,
		cfg:     cfg,
	}, nil
}

func (p *Producer) Shutdown() error {
	p.log.Info("closing scheduler channel...")
	err := p.channel.Close()
	if err != nil {
		return err
	}
	p.log.Info("closing scheduler connection...")
	err = p.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) Publish(ctx context.Context, body []byte) error {
	if p.channel == nil {
		return ErrSchedulerRabbitNilChannel
	}

	p.log.Info("publishing...")
	err := p.channel.PublishWithContext(
		ctx,
		p.cfg.ExchangeName,
		p.cfg.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "utf8",
			DeliveryMode:    uint8(p.cfg.DeliveryMode),
			Timestamp:       time.Now(),
			Body:            body,
		},
	)
	if err != nil {
		return err
	}

	p.log.Info("message is published...")

	return nil
}
