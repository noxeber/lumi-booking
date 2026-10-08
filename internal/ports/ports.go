package ports

import (
	"time"
	"uuid"

	"github.com/noxeber/lumi-booking/internal/domain/appointment"
	"github.com/noxeber/lumi-booking/internal/domain/schedule"
	"github.com/noxeber/lumi-booking/internal/domain/service"
)

type ServiceRepo interface {
	GetByID(id uuid.UUID) (service.Service, error)
}

type ScheduleRepo interface {
	GetByDateAndMaster(date time.Time, masterID uuid.UUID) (schedule.Schedule, error)
}

type AppointmentRepo interface {
	GetByDateAndMaster(date time.Time, masterID uuid.UUID) ([]appointment.Appointment, error)
}
