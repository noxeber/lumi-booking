package appointment

import (
	"errors"
	"time"
	"uuid"
)

var ErrAlreadyExist = errors.New("id is exist in slice ids of appointment")

type Appointment struct {
	id         uuid.UUID
	userID     uuid.UUID
	masterID   uuid.UUID
	serviceIDs []uuid.UUID
	timeStart  time.Time
	createdAt  time.Time
}

func NewAppointment(userID, masterID uuid.UUID, serviceIDs []uuid.UUID, timeStart time.Time) (*Appointment, error) {
	return &Appointment{id: uuid.NewV7(), userID: userID, masterID: masterID, serviceIDs: serviceIDs, timeStart: timeStart, createdAt: time.Now()}, nil
}

func (a *Appointment) AddIDService(id uuid.UUID) error {
	for _, exist := range a.serviceIDs {
		if id == exist {
			return ErrAlreadyExist
		}
	}
	a.serviceIDs = append(a.serviceIDs, id)
	return nil
}
