package rabbitmq

type Message struct {
	EventID     string `json:"event_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        int64  `json:"date"`
	UserID      string `json:"user_id"`
}

type Notification struct {
	Message Message
	Err     error
}
