package application

import (
	"errors"
	"sort"
	"time"

	"uuid"

	"github.com/noxeber/lumi-booking/internal/domain/appointment"
	"github.com/noxeber/lumi-booking/internal/domain/schedule"
	"github.com/noxeber/lumi-booking/internal/ports"
)

var ErrIsDayOff = errors.New("schedule is day off")

type TimeSlot struct {
	StartTime time.Time
	EndTime   time.Time
}

type PeriodSlots struct {
	Name  string
	Slots []time.Time
}

type SlotService struct {
	serviceRepo     ports.ServiceRepo
	scheduleRepo    ports.ScheduleRepo
	appointmentRepo ports.AppointmentRepo
}

func NewSlotService(serviceRepo ports.ServiceRepo, scheduleRepo ports.ScheduleRepo, appointmentRepo ports.AppointmentRepo) *SlotService {
	return &SlotService{
		serviceRepo:     serviceRepo,
		scheduleRepo:    scheduleRepo,
		appointmentRepo: appointmentRepo,
	}
}

var ErrNoServicesSelected = errors.New("no services selected")

func (s *SlotService) GetAvailableSlots(date time.Time, masterID uuid.UUID, serviceIDs []uuid.UUID) ([]PeriodSlots, error) {
	if len(serviceIDs) == 0 {
		return nil, ErrNoServicesSelected
	}

	var dur time.Duration
	var maxBuffer time.Duration

	for _, serviceID := range serviceIDs {
		srv, err := s.serviceRepo.GetByID(serviceID)
		if err != nil {
			return nil, err
		}
		dur += srv.Duration()
		if srv.BufferTime() > maxBuffer {
			maxBuffer = srv.BufferTime()
		}
	}

	totalRequiredTime := dur + maxBuffer

	sched, err := s.scheduleRepo.GetByDateAndMaster(date, masterID)
	if err != nil {
		return nil, err
	}

	if sched.IsDayOff() {
		return nil, ErrIsDayOff
	}

	appointments, err := s.appointmentRepo.GetByDateAndMaster(date, masterID)
	if err != nil {
		return nil, err
	}

	windows := getFreeWindows(sched, appointments)
	filteredWindows := filterWindows(windows, totalRequiredTime)
	steps := generateSteps(filteredWindows, totalRequiredTime, 15*time.Minute)

	return groupIntoPeriods(steps), nil
}

func getFreeWindows(sched schedule.Schedule, appointments []appointment.Appointment) []TimeSlot {
	type busyBlock struct {
		start time.Time
		end   time.Time
	}

	var busyBlocks []busyBlock

	// Add breaks
	for _, b := range sched.Breaks() {
		busyBlocks = append(busyBlocks, busyBlock{start: b.StartTime(), end: b.EndTime()})
	}

	// Add appointments (including their buffer times)
	for _, ap := range appointments {
		if ap.Status() == appointment.StatusActive || ap.Status() == appointment.StatusCompleted {
			busyBlocks = append(busyBlocks, busyBlock{
				start: ap.TimeStart(),
				end:   ap.TimeEnd().Add(ap.BufferTime()),
			})
		}
	}

	// Sort blocks by start time
	sort.Slice(busyBlocks, func(i, j int) bool {
		return busyBlocks[i].start.Before(busyBlocks[j].start)
	})

	var freeWindows []TimeSlot
	currentStart := sched.StartTime()

	for _, block := range busyBlocks {
		if block.start.After(currentStart) {
			freeWindows = append(freeWindows, TimeSlot{StartTime: currentStart, EndTime: block.start})
		}
		if block.end.After(currentStart) {
			currentStart = block.end
		}
	}

	if currentStart.Before(sched.EndTime()) {
		freeWindows = append(freeWindows, TimeSlot{StartTime: currentStart, EndTime: sched.EndTime()})
	}

	return freeWindows
}

func filterWindows(windows []TimeSlot, requiredTime time.Duration) []TimeSlot {
	var filtered []TimeSlot
	for _, w := range windows {
		if w.EndTime.Sub(w.StartTime) >= requiredTime {
			filtered = append(filtered, w)
		}
	}
	return filtered
}

func generateSteps(windows []TimeSlot, requiredTime time.Duration, step time.Duration) []time.Time {
	var steps []time.Time
	for _, w := range windows {
		current := w.StartTime
		for current.Add(requiredTime).Before(w.EndTime) || current.Add(requiredTime).Equal(w.EndTime) {
			steps = append(steps, current)
			current = current.Add(step)
		}
	}
	return steps
}

func groupIntoPeriods(steps []time.Time) []PeriodSlots {
	morning := PeriodSlots{Name: "Утро"}
	day := PeriodSlots{Name: "День"}
	evening := PeriodSlots{Name: "Вечер"}

	for _, step := range steps {
		hour := step.Hour()
		if hour < 12 {
			morning.Slots = append(morning.Slots, step)
		} else if hour < 17 {
			day.Slots = append(day.Slots, step)
		} else {
			evening.Slots = append(evening.Slots, step)
		}
	}

	var result []PeriodSlots
	if len(morning.Slots) > 0 {
		result = append(result, morning)
	}
	if len(day.Slots) > 0 {
		result = append(result, day)
	}
	if len(evening.Slots) > 0 {
		result = append(result, evening)
	}

	return result
}
