package service

//go:generate mockgen -source=service.go -destination=mock/mock.go service

import (
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/logger"
	"github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage"
)

type Service struct {
	storage storage.Storage
	logger  logger.Logger
}

type Services interface {
	storage.Storage
	logger.Logger
}

func NewService(storage storage.Storage, logger logger.Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}
