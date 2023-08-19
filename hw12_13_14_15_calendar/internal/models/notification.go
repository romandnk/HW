package models

import "time"

type Notification struct {
	EventID  string
	Title    string
	Date     time.Time
	UserID   int
	Interval time.Duration
}
