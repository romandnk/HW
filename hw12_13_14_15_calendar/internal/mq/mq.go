package mq

import "context"

type Producer interface {
	Publish(ctx context.Context, body []byte) error
}

type Consumer interface {
	Consume() error
}
