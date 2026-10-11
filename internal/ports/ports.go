package ports

import (
	"github.com/google/uuid"
	"time"

	"github.com/noxeber/lumi-booking/internal/domain/appointment"
	"github.com/noxeber/lumi-booking/internal/domain/master"
	"github.com/noxeber/lumi-booking/internal/domain/schedule"
	"github.com/noxeber/lumi-booking/internal/domain/service"
	"github.com/noxeber/lumi-booking/internal/domain/user"
)

type ServiceRepo interface {
	GetByID(id uuid.UUID) (*service.Service, error)
	Save(srv *service.Service) error
}

type ScheduleRepo interface {
	GetByDateAndMaster(date time.Time, masterID uuid.UUID) (*schedule.Schedule, error)
	Save(sched *schedule.Schedule) error
}

type AppointmentRepo interface {
	GetByDateAndMaster(date time.Time, masterID uuid.UUID) ([]appointment.Appointment, error)
	GetByID(id uuid.UUID) (*appointment.Appointment, error)
	Save(appt *appointment.Appointment) error
}

type MasterRepo interface {
	GetByID(id uuid.UUID) (*master.Master, error)
	Save(m *master.Master) error
}

type UserRepo interface {
	GetByID(id uuid.UUID) (*user.User, error)
	Save(u *user.User) error
}
