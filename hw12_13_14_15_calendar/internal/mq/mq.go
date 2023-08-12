package mq

import "context"

type OpenerCloserChannel interface {
	OpenChannel() error
	CloseChannel() error
}

type CloserConn interface {
	CloseConn() error
}

type Scheduler interface {
	OpenerCloserChannel
	CloserConn
	Publish(ctx context.Context, body []byte) error
}

type Sender interface {
	OpenerCloserChannel
	CloserConn
	Consume(ctx context.Context) error
}
