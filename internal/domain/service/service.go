package service

import (
	"errors"
	"time"
	"unicode/utf8"
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
	bufferTime  time.Duration
}

func CreateService(name string, price decimal.Decimal, duration time.Duration, bufferTime time.Duration, description string) (*Service, error) {
	if utf8.RuneCountInString(name) < 4 {
		return nil, ErrInvalidName
	}
	if price.IsNegative() {
		return nil, ErrNegativePrice
	}
	return &Service{id: uuid.New(), name: name, price: price, duration: duration, bufferTime: bufferTime, description: description}, nil
}

func RestoreService(id uuid.UUID, name string, price decimal.Decimal, duration time.Duration, bufferTime time.Duration, description string) *Service {
	return &Service{id: id, name: name, price: price, duration: duration, bufferTime: bufferTime, description: description}
}

func (s Service) Duration() time.Duration {
	return s.duration
}

func (s Service) BufferTime() time.Duration {
	return s.bufferTime
}
