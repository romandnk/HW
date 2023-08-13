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
	cfg     config.SchedulerRabbit
}

func NewScheduler(cfg config.SchedulerRabbit, log logger.Logger) (*Scheduler, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conf := amqp.Config{
		Heartbeat: cfg.Heartbeat,
	}
	conn, err := amqp.DialConfig(url, conf)
	if err != nil {
		return nil, err
	}

	go func() {
		conn.NotifyClose(make(chan *amqp.Error))
	}()

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
		s.cfg.Durable,
		s.cfg.AutoDelete,
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

	s.log.Info("publishing...")
	err := s.channel.PublishWithContext(
		ctx,
		s.cfg.ExchangeName,
		s.cfg.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "utf8",
			DeliveryMode:    s.cfg.DeliveryMode,
			Timestamp:       time.Now(),
			Body:            body,
		},
	)
	if err != nil {
		s.log.Error("error publishing", slog.String("error", err.Error()))
		return err
	}

	return nil
}
