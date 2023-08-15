package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"golang.org/x/exp/slog"
)

var ErrSenderRabbitNilChannel = errors.New("rabbit sender: channel is nil")

type SenderConfig struct {
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

type Sender struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	log     logger.Logger
	cfg     SenderConfig
	done    chan struct{}
}

func NewSender(cfg SenderConfig, log logger.Logger) (*Sender, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conf := amqp.Config{
		Heartbeat: cfg.Heartbeat,
	}

	log.Info("connecting sender...")
	conn, err := amqp.DialConfig(url, conf)
	if err != nil {
		return nil, err
	}

	log.Info("opening sender channel...")
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	log.Info("declaring sender exchange...")
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

	log.Info("declaring sender queue...")
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

	log.Info("binding sender exchange...")
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

	return &Sender{
		conn:    conn,
		channel: ch,
		log:     log,
		cfg:     cfg,
	}, nil
}

func (s *Sender) Consume() error {
	if s.channel == nil {
		return ErrSenderRabbitNilChannel
	}

	s.log.Info("starting consuming...")
	deliveries, err := s.channel.Consume(
		s.cfg.QueueName,
		s.cfg.Tag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.log.Error("error while starting consuming", slog.String("error", err.Error()))
		return err
	}

	notifications := make(chan Notification)

	go handle(deliveries, notifications, s.done)

	for {
		select {
		case <-s.done:
			close(notifications)
			return nil
		case notification, ok := <-notifications:
			if !ok {
				return nil
			}

			if notification.Err != nil {
				s.log.Error("error receiving notification", slog.String("error", notification.Err.Error()))
				continue
			}
			s.log.Info("notification is received", slog.Any("notification", notification.Message))
		}
	}
}

func (s *Sender) Shutdown() error {
	s.log.Info("closing deliveries...")
	err := s.channel.Cancel(s.cfg.Tag, true)
	if err != nil {
		return err
	}
	s.log.Info("closing sender connection...")
	err = s.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func handle(deliveries <-chan amqp.Delivery, notifications chan Notification, done chan struct{}) {
	for delivery := range deliveries {
		var msg Message
		err := json.Unmarshal(delivery.Body, &msg)

		if err == nil {
			delivery.Ack(false)
		}

		notification := Notification{
			Message: msg,
			Err:     err,
		}

		notifications <- notification
	}

	done <- struct{}{}
}
