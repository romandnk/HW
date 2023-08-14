package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/cmd/config"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"golang.org/x/exp/slog"
)

var ErrSenderRabbitNilChannel = errors.New("rabbit sender: channel is nil")

type Sender struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	deliveries <-chan amqp.Delivery
	log        logger.Logger
	cfg        config.RabbitConfig
}

func NewSender(cfg config.RabbitConfig, log logger.Logger) (*Sender, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conf := amqp.Config{
		Heartbeat: cfg.Heartbeat,
	}
	conn, err := amqp.DialConfig(url, conf)
	if err != nil {
		return nil, err
	}

	return &Sender{
		conn: conn,
		log:  log,
		cfg:  cfg,
	}, nil
}

func (s *Sender) CloseConn() error {
	s.log.Info("closing sender connection...")
	err := s.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) OpenChannel() error {
	s.log.Info("opening sender channel...")
	ch, err := s.conn.Channel()
	if err != nil {
		s.log.Error("error opening sender chanel", slog.String("error", err.Error()))
		return err
	}
	s.channel = ch

	s.log.Info("declaring sender exchange...")
	err = s.channel.ExchangeDeclare(
		s.cfg.ExchangeName,
		s.cfg.ExchangeType,
		s.cfg.DurableExchange,
		s.cfg.AutoDeleteExchange,
		false,
		false,
		nil,
	)
	if err != nil {
		s.log.Error("error declaring sender exchange", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *Sender) CloseChannel() error {
	s.log.Info("closing sender channel...")
	err := s.channel.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) Consume() error {
	if s.channel == nil {
		return ErrSenderRabbitNilChannel
	}

	s.log.Info("declaring sender queue...")
	queue, err := s.channel.QueueDeclare(
		s.cfg.QueueName,
		s.cfg.DurableQueue,
		s.cfg.AutoDeleteQueue,
		false,
		false,
		nil,
	)
	if err != nil {
		s.log.Error("error declaring sender queue", slog.String("error", err.Error()))
		return err
	}

	s.log.Info("binding sender exchange...")
	err = s.channel.QueueBind(
		queue.Name,
		s.cfg.RoutingKey,
		s.cfg.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		s.log.Error("error binding sender exchange", slog.String("error", err.Error()))
		return err
	}

	s.log.Info("starting consuming...")
	deliveries, err := s.channel.Consume(
		queue.Name,
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

	s.deliveries = deliveries

	return nil
}

func (s *Sender) Handle(ctx context.Context) <-chan Notification {
	messages := make(chan Notification)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case delivery := <-s.deliveries:
				var msg Message
				var notif Notification
				err := json.Unmarshal(delivery.Body, &msg)
				if err != nil {
					notif = Notification{
						Message: Message{},
						Err:     err,
					}

					messages <- notif
					return
				}

				if err := delivery.Ack(false); err != nil {
					notif = Notification{
						Message: Message{},
						Err:     err,
					}

					messages <- notif
					return
				}

				notif = Notification{
					Message: msg,
					Err:     nil,
				}

				messages <- notif
			}
		}
	}()

	return messages
}
