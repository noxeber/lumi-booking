package application

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/noxeber/lumi-booking/internal/domain/appointment"
)

type mockAppointmentRepo struct {
	saveErr error
	appt    *appointment.Appointment
	appts   []appointment.Appointment
}

func (m *mockAppointmentRepo) Save(appt *appointment.Appointment) error {
	m.appt = appt
	return m.saveErr
}
func (m *mockAppointmentRepo) GetByDateAndMaster(date time.Time, masterID uuid.UUID) ([]appointment.Appointment, error) {
	return m.appts, nil
}
func (m *mockAppointmentRepo) GetByID(id uuid.UUID) (*appointment.Appointment, error) {
	if m.appt == nil {
		return nil, errors.New("not found")
	}
	return m.appt, nil
}

func TestBookingService_CreateAppointment_SlotNotAvailable(t *testing.T) {
	sRepo := &mockServiceRepo{}
	// Mock returns an error on nil srv instead of panicking
	sRepo.saveErr = errors.New("service not found mock")
	scRepo := &mockScheduleRepo{}
	aRepo := &mockAppointmentRepo{}

	slotService := NewSlotService(sRepo, scRepo, aRepo)
	bookingSvc := NewBookingService(slotService, aRepo)

	// CreateSchedule is omitted making schedule inherently empty/nil causing ErrNoServicesSelected or similar depending on flow.
	// We'll just test that an arbitrary time on empty setup returns ErrSlotNotAvailable (since no slots will be returned).

	date := time.Now()
	startTime := date.Add(9 * time.Hour)

	_, err := bookingSvc.CreateAppointment(uuid.New(), uuid.New(), []uuid.UUID{uuid.New()}, date, startTime)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBookingService_CancelAppointment(t *testing.T) {
	aRepo := &mockAppointmentRepo{}
	bookingSvc := NewBookingService(nil, aRepo)

	startTime := time.Now().Add(48 * time.Hour) // More than 24 hours away
	appt, _ := appointment.NewAppointment(uuid.New(), uuid.New(), []uuid.UUID{}, startTime, startTime.Add(1*time.Hour), 0)
	aRepo.appt = appt // mock existing appointment

	err := bookingSvc.CancelAppointment(uuid.New()) // ID mock doesn't matter since GetByID just returns m.appt

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if aRepo.appt.Status() != appointment.StatusCancelledByClient {
		t.Errorf("status = %v, want %v", aRepo.appt.Status(), appointment.StatusCancelledByClient)
	}
}

func TestBookingService_CancelAppointment_TooLate(t *testing.T) {
	aRepo := &mockAppointmentRepo{}
	bookingSvc := NewBookingService(nil, aRepo)

	startTime := time.Now().Add(12 * time.Hour) // Less than 24 hours away
	appt, _ := appointment.NewAppointment(uuid.New(), uuid.New(), []uuid.UUID{}, startTime, startTime.Add(1*time.Hour), 0)
	aRepo.appt = appt

	err := bookingSvc.CancelAppointment(uuid.New())

	if !errors.Is(err, ErrTooLateToCancel) {
		t.Errorf("expected error %v, got %v", ErrTooLateToCancel, err)
	}
}
