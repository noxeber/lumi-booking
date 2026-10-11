package application

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/noxeber/lumi-booking/internal/domain/appointment"
	"github.com/noxeber/lumi-booking/internal/ports"
)

var (
	ErrSlotNotAvailable    = errors.New("requested time slot is not available")
	ErrAppointmentNotFound = errors.New("appointment not found")
	ErrTooLateToCancel     = errors.New("cancellation must be done at least 24 hours before start time")
)

type BookingService struct {
	slotService     *SlotService
	appointmentRepo ports.AppointmentRepo
}

func NewBookingService(slotService *SlotService, appointmentRepo ports.AppointmentRepo) *BookingService {
	return &BookingService{
		slotService:     slotService,
		appointmentRepo: appointmentRepo,
	}
}

func (s *BookingService) CreateAppointment(userID, masterID uuid.UUID, serviceIDs []uuid.UUID, date, startTime time.Time) (*appointment.Appointment, error) {
	// First verify the requested slot is actually available
	availablePeriods, err := s.slotService.GetAvailableSlots(date, masterID, serviceIDs)
	if err != nil {
		return nil, err
	}

	isAvailable := false
	for _, period := range availablePeriods {
		for _, slot := range period.Slots {
			if slot.Equal(startTime) {
				isAvailable = true
				break
			}
		}
		if isAvailable {
			break
		}
	}

	if !isAvailable {
		return nil, ErrSlotNotAvailable
	}

	// Calculate total duration and max buffer
	var dur time.Duration
	var maxBuffer time.Duration
	for _, serviceID := range serviceIDs {
		srv, err := s.slotService.serviceRepo.GetByID(serviceID)
		if err != nil {
			return nil, err
		}
		dur += srv.Duration()
		if srv.BufferTime() > maxBuffer {
			maxBuffer = srv.BufferTime()
		}
	}

	endTime := startTime.Add(dur)
	appt, err := appointment.NewAppointment(userID, masterID, serviceIDs, startTime, endTime, maxBuffer)
	if err != nil {
		return nil, err
	}

	if err := s.appointmentRepo.Save(appt); err != nil {
		return nil, err
	}

	return appt, nil
}

func (s *BookingService) CancelAppointment(apptID uuid.UUID) error {
	appt, err := s.appointmentRepo.GetByID(apptID)
	if err != nil {
		return err
	}
	if appt == nil {
		return ErrAppointmentNotFound
	}

	// Cancellation is only allowed if more than 24 hours left before StartTime (as per TZ)
	timeUntilStart := time.Until(appt.TimeStart())
	if timeUntilStart < 24*time.Hour {
		return ErrTooLateToCancel
	}

	// Note: We would ideally add a SetStatus method to domain model for cleaner manipulation.
	// Since TZ requires domain models to encapsulate logic we will simulate saving the state change.
	// We'll update the domain model `Appointment` slightly in bash later.
	appt.CancelByClient()

	return s.appointmentRepo.Save(appt)
}
