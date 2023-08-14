package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/cmd/config"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"golang.org/x/exp/slog"
)

var ErrSchedulerRabbitNilChannel = errors.New("rabbit scheduler: channel is nil")

type Scheduler struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	log     logger.Logger
	cfg     config.RabbitSchedulerConfig
}

func NewScheduler(cfg config.RabbitSchedulerConfig, log logger.Logger) (*Scheduler, error) {
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

	return &Scheduler{
		conn: conn,
		log:  log,
		cfg:  cfg,
	}, nil
}

func (s *Scheduler) CloseConn() error {
	s.log.Info("closing scheduler connection...")
	err := s.conn.Close()
	if err != nil {
		s.log.Error("error closing scheduler connection", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *Scheduler) OpenChannel() error {
	s.log.Info("opening scheduler channel...")
	ch, err := s.conn.Channel()
	if err != nil {
		s.log.Error("error opening scheduler chanel", slog.String("error", err.Error()))
		return err
	}
	s.channel = ch

	s.log.Info("declaring scheduler exchange...")
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
		s.log.Error("error declaring scheduler exchange", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *Scheduler) CloseChannel() error {
	s.log.Info("closing scheduler channel...")
	err := s.channel.Close()
	if err != nil {
		s.log.Error("error closing scheduler channel", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *Scheduler) Publish(ctx context.Context, body []byte) error {
	if s.channel == nil {
		return ErrSchedulerRabbitNilChannel
	}

	queue, err := s.channel.QueueDeclare(
		s.cfg.QueueName,
		s.cfg.DurableQueue,
		s.cfg.AutoDeleteQueue,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = s.channel.QueueBind(
		queue.Name,
		s.cfg.RoutingKey,
		s.cfg.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	s.log.Info("publishing...")
	err = s.channel.PublishWithContext(
		ctx,
		s.cfg.ExchangeName,
		s.cfg.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "utf8",
			DeliveryMode:    uint8(s.cfg.DeliveryMode),
			Timestamp:       time.Now(),
			Body:            body,
		},
	)
	if err != nil {
		s.log.Error("error publishing", slog.String("error", err.Error()))
		return err
	}

	s.log.Info("message is published...")

	return nil
}
