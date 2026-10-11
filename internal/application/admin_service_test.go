package application

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/noxeber/lumi-booking/internal/domain/master"
	"github.com/noxeber/lumi-booking/internal/domain/schedule"
	"github.com/noxeber/lumi-booking/internal/domain/service"
)

type mockServiceRepo struct {
	saveErr error
	srv     *service.Service
}

func (m *mockServiceRepo) Save(s *service.Service) error {
	m.srv = s
	return m.saveErr
}
func (m *mockServiceRepo) GetByID(id uuid.UUID) (*service.Service, error) {
	if m.srv == nil {
		return nil, errors.New("not found")
	}
	return m.srv, nil
}

type mockMasterRepo struct {
	saveErr error
	m       *master.Master
}

func (m *mockMasterRepo) Save(master *master.Master) error {
	m.m = master
	return m.saveErr
}
func (m *mockMasterRepo) GetByID(id uuid.UUID) (*master.Master, error) {
	if m.m == nil {
		return nil, errors.New("not found")
	}
	return m.m, nil
}

type mockScheduleRepo struct {
	saveErr error
	s       *schedule.Schedule
}

func (m *mockScheduleRepo) Save(sched *schedule.Schedule) error {
	m.s = sched
	return m.saveErr
}
func (m *mockScheduleRepo) GetByDateAndMaster(date time.Time, masterID uuid.UUID) (*schedule.Schedule, error) {
	if m.s == nil {
		return nil, errors.New("not found")
	}
	return m.s, nil
}

func TestAdminService_CreateService(t *testing.T) {
	repo := &mockServiceRepo{}
	svc := NewAdminService(repo, nil, nil)

	srv, err := svc.CreateService("Haircut", decimal.NewFromInt(100), time.Hour, 15*time.Minute, "desc")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if srv == nil {
		t.Fatalf("expected service, got nil")
	}
	if repo.srv == nil {
		t.Error("service was not saved to repo")
	}
}

func TestAdminService_CreateMaster(t *testing.T) {
	repo := &mockMasterRepo{}
	svc := NewAdminService(nil, repo, nil)

	m, err := svc.CreateMaster("John Doe")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatalf("expected master, got nil")
	}
	if repo.m == nil {
		t.Error("master was not saved to repo")
	}
}

func TestAdminService_AssignServiceToMaster(t *testing.T) {
	repo := &mockMasterRepo{}
	svc := NewAdminService(nil, repo, nil)

	_, _ = svc.CreateMaster("John")
	serviceID := uuid.New()

	err := svc.AssignServiceToMaster(uuid.New(), serviceID) // Note: uuid.New() vs m.ID won't matter for this mock as GetByID returns the only master

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Note: checking master services count properly requires inspecting the unexported field, but error nil means AddService passed
}

func TestAdminService_CreateSchedule(t *testing.T) {
	repo := &mockScheduleRepo{}
	svc := NewAdminService(nil, nil, repo)

	date := time.Now()
	start := date.Add(9 * time.Hour)
	end := date.Add(17 * time.Hour)

	sched, err := svc.CreateSchedule(uuid.New(), date, start, end, false)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sched == nil {
		t.Fatalf("expected schedule, got nil")
	}
	if repo.s == nil {
		t.Error("schedule was not saved to repo")
	}
}

func TestAdminService_AddBreakToSchedule(t *testing.T) {
	repo := &mockScheduleRepo{}
	svc := NewAdminService(nil, nil, repo)

	date := time.Now()
	start := date.Add(9 * time.Hour)
	end := date.Add(17 * time.Hour)

	svc.CreateSchedule(uuid.New(), date, start, end, false)

	breakStart := date.Add(13 * time.Hour)
	breakEnd := date.Add(14 * time.Hour)

	err := svc.AddBreakToSchedule(uuid.New(), date, breakStart, breakEnd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
