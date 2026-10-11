package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/noxeber/lumi-booking/internal/domain/master"
	"github.com/noxeber/lumi-booking/internal/domain/schedule"
	"github.com/noxeber/lumi-booking/internal/domain/service"
	"github.com/noxeber/lumi-booking/internal/ports"
)

type AdminService struct {
	serviceRepo  ports.ServiceRepo
	masterRepo   ports.MasterRepo
	scheduleRepo ports.ScheduleRepo
}

func NewAdminService(serviceRepo ports.ServiceRepo, masterRepo ports.MasterRepo, scheduleRepo ports.ScheduleRepo) *AdminService {
	return &AdminService{
		serviceRepo:  serviceRepo,
		masterRepo:   masterRepo,
		scheduleRepo: scheduleRepo,
	}
}

func (s *AdminService) CreateService(name string, price decimal.Decimal, duration, bufferTime time.Duration, description string) (*service.Service, error) {
	srv, err := service.CreateService(name, price, duration, bufferTime, description)
	if err != nil {
		return nil, err
	}

	if err := s.serviceRepo.Save(srv); err != nil {
		return nil, err
	}
	return srv, nil
}

func (s *AdminService) CreateMaster(name string) (*master.Master, error) {
	m, err := master.CreateMaster(name)
	if err != nil {
		return nil, err
	}

	if err := s.masterRepo.Save(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *AdminService) AssignServiceToMaster(masterID, serviceID uuid.UUID) error {
	m, err := s.masterRepo.GetByID(masterID)
	if err != nil {
		return err
	}

	if err := m.AddService(serviceID); err != nil {
		return err
	}

	return s.masterRepo.Save(m)
}

func (s *AdminService) CreateSchedule(masterID uuid.UUID, date, startTime, endTime time.Time, isDayOff bool) (*schedule.Schedule, error) {
	sched, err := schedule.NewSchedule(masterID, date, startTime, endTime, isDayOff)
	if err != nil {
		return nil, err
	}

	if err := s.scheduleRepo.Save(sched); err != nil {
		return nil, err
	}
	return sched, nil
}

func (s *AdminService) AddBreakToSchedule(masterID uuid.UUID, date, breakStart, breakEnd time.Time) error {
	sched, err := s.scheduleRepo.GetByDateAndMaster(date, masterID)
	if err != nil {
		return err
	}

	b, err := schedule.NewBreak(breakStart, breakEnd)
	if err != nil {
		return err
	}

	if err := sched.AddBreak(b); err != nil {
		return err
	}

	return s.scheduleRepo.Save(sched)
}
