package application

import (
	"errors"
	"sort"
	"time"
	"uuid"

	"github.com/noxeber/lumi-booking/internal/domain/appointment"
	schedule "github.com/noxeber/lumi-booking/internal/domain/workday"
	"github.com/noxeber/lumi-booking/internal/ports"
)

var ErrIsDayOff = errors.New("workday is dayoff")

type TimeSlot struct {
	startTime time.Time
	endTime   time.Time
}

type SlotService struct {
	serviceRepo     ports.ServiceRepo
	workdayRepo     ports.WorkdayRepo
	appointmentRepo ports.AppointmentRepo
}

func (s *SlotService) GetAvaliableSlots(date time.Time, masterID uuid.UUID, serviceIDs []uuid.UUID) ([]TimeSlot, error) {
	var dur time.Duration

	for _, serviceID := range serviceIDs {
		service, err := s.serviceRepo.GetByID(serviceID)
		if err != nil {
			return nil, err
		}
		dur += service.Duration()
	}

	workDay, err := s.workdayRepo.GetByDateAndMaster(date, masterID)
	if err != nil {
		return nil, err
	}

	for _, action := range workDay.Actions() {
		if typeAction := action.TypeAction(); typeAction == schedule.ActionTypeDayOff {
			return nil, ErrIsDayOff
		}
	}

	appointments, err := s.appointmentRepo.GetByDateAndMaster(date, masterID)
	if err != nil {
		return nil, err
	}

	return calculateSlots(appointments, workDay, dur), nil
}

func calculateSlots(appointments []appointment.Appointment, workday schedule.WorkDay, totalDuration time.Duration) []TimeSlot {
	if totalDuration <= 0 {
		return nil
	}

	freeIntervals := workday.FreeIntervals()
	var availableSlots []TimeSlot

	for _, free := range freeIntervals {
		currentStart := free.StartTime()

		var overlapping []appointment.Appointment
		for _, ap := range appointments {
			if ap.TimeStart().Before(free.EndTime()) && ap.TimeEnd().After(free.StartTime()) {
				overlapping = append(overlapping, ap)
			}
		}

		sort.Slice(overlapping, func(i, j int) bool {
			return overlapping[i].TimeStart().Before(overlapping[j].TimeStart())
		})

		for _, ap := range overlapping {
			if ap.TimeStart().After(currentStart) {
				windowDuration := ap.TimeStart().Sub(currentStart)
				if windowDuration >= totalDuration {
					availableSlots = append(availableSlots, TimeSlot{
						startTime: currentStart,
						endTime:   ap.TimeStart(),
					})
				}
			}
			if ap.TimeEnd().After(currentStart) {
				currentStart = ap.TimeEnd()
			}
		}

		if free.EndTime().After(currentStart) {
			windowDuration := free.EndTime().Sub(currentStart)
			if windowDuration >= totalDuration {
				availableSlots = append(availableSlots, TimeSlot{
					startTime: currentStart,
					endTime:   free.EndTime(),
				})
			}
		}
	}

	return availableSlots
}
