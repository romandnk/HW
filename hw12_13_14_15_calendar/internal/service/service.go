package service

import (
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage"
)

type Service struct {
	storage Storage
	logger  Logger
}

type Storage interface {
	storage.StoreEvent
}

type Logger interface {
	WriteLogInFile(path string) error
}

func NewService(storage Storage, logger Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}
