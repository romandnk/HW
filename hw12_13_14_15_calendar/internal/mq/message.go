package mq

import "time"

type Message struct {
	EventID string    `json:"event_id"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	UserID  int       `json:"user_id"`
}

type Notification struct {
	Message Message
	Err     error
}
