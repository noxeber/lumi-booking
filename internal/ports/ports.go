package ports

import (
	"time"
	"uuid"

	"github.com/noxeber/lumi-booking/internal/domain/appointment"
	"github.com/noxeber/lumi-booking/internal/domain/service"
	schedule "github.com/noxeber/lumi-booking/internal/domain/workday"
)

type ServiceRepo interface {
	GetByID(id uuid.UUID) (service.Service, error)
}

type WorkdayRepo interface {
	GetByDateAndMaster(date time.Time, masteID uuid.UUID) (schedule.WorkDay, error)
}

type AppointmentRepo interface {
	GetByDateAndMaster(date time.Time, masterID uuid.UUID) ([]appointment.Appointment, error)
}
