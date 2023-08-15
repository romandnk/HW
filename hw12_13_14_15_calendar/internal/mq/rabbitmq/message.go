package rabbitmq

import "time"

type Message struct {
	EventID     string    `json:"event_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	UserID      string    `json:"user_id"`
}

type Notification struct {
	Message Message
	Err     error
}
