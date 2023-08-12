package config

import "time"

type SchedulerRabbit struct {
	Username     string
	Password     string
	Host         string
	Port         int
	Heartbeat    time.Duration
	ExchangeName string
	ExchangeType string
	Durable      bool
	AutoDelete   bool
	RoutingKey   string
	DeliveryMode uint8
}
