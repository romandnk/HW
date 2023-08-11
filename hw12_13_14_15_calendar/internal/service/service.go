package service

//go:generate mockgen -source=service.go -destination=mock/mock.go service

import (
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage"
)

type Service struct {
	event storage.Event
}

type Services interface {
	storage.Event
}

func NewService(event storage.Event) *Service {
	return &Service{
		event: event,
	}
}
