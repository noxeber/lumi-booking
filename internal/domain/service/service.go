package service

import (
	"errors"
	"time"
	"uuid"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidName   = errors.New("service name is empty or too short")
	ErrNegativePrice = errors.New("service price must be positive")
)

type Service struct {
	id          uuid.UUID
	name        string
	description string
	price       decimal.Decimal
	duration    time.Duration
}

func CreateService(name string, price decimal.Decimal, duration time.Duration, description string) (*Service, error) {
	if len(name) < 4 {
		return nil, ErrInvalidName
	}
	if price.IsNegative() {
		return nil, ErrNegativePrice
	}
	return &Service{id: uuid.NewV4(), name: name, price: price, duration: duration, description: description}, nil
}

func RestoreService(id uuid.UUID, name string, price decimal.Decimal, duration time.Duration, description string) *Service {
	return &Service{id: id, name: name, price: price, duration: duration, description: description}
}

func (s Service) Duration() time.Duration {
	return s.duration
}
