package mq

import "context"

type NotificationProducer interface {
	Publish(ctx context.Context, body []byte) error
	Shutdown() error
}

type Producer struct {
	NotificationProducer
}

func NewProducer(notificationProducer NotificationProducer) *Producer {
	return &Producer{NotificationProducer: notificationProducer}
}

type NotificationConsumer interface {
	Consume() (<-chan Notification, error)
	Shutdown() error
}

type Consumer struct {
	NotificationConsumer
}

func NewConsumer(notificationConsumer NotificationConsumer) *Consumer {
	return &Consumer{NotificationConsumer: notificationConsumer}
}
