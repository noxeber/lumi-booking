package appointment

import (
	"errors"
	"time"
	"uuid"
)

var ErrAlreadyExist = errors.New("id is exist in slice ids of appointment")

type AppointmentStatus string

const (
	StatusActive            AppointmentStatus = "Active"
	StatusCompleted         AppointmentStatus = "Completed"
	StatusCancelledByClient AppointmentStatus = "CancelledByClient"
	StatusCancelledByAdmin  AppointmentStatus = "CancelledByAdmin"
)

type Appointment struct {
	id         uuid.UUID
	userID     uuid.UUID
	masterID   uuid.UUID
	serviceIDs []uuid.UUID
	timeStart  time.Time
	timeEnd    time.Time
	bufferTime time.Duration
	createdAt  time.Time
	status     AppointmentStatus
}

func NewAppointment(userID, masterID uuid.UUID, serviceIDs []uuid.UUID, timeStart, timeEnd time.Time, bufferTime time.Duration) (*Appointment, error) {
	return &Appointment{
		id:         uuid.New(),
		userID:     userID,
		masterID:   masterID,
		serviceIDs: serviceIDs,
		timeStart:  timeStart,
		timeEnd:    timeEnd,
		bufferTime: bufferTime,
		createdAt:  time.Now(),
		status:     StatusActive,
	}, nil
}

func (a *Appointment) BufferTime() time.Duration {
	return a.bufferTime
}

func (a *Appointment) Status() AppointmentStatus {
	return a.status
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

func (a *Appointment) TimeStart() time.Time {
	return a.timeStart
}

func (a *Appointment) TimeEnd() time.Time {
	return a.timeEnd
}
