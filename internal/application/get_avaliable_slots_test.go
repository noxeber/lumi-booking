package application

import (
	"testing"
	"time"

	"uuid"

	"github.com/noxeber/lumi-booking/internal/domain/appointment"
	"github.com/noxeber/lumi-booking/internal/domain/schedule"
)

func TestGetFreeWindows(t *testing.T) {
	masterID := uuid.New()
	date := time.Now().Truncate(24 * time.Hour)
	start := date.Add(9 * time.Hour) // 09:00
	end := date.Add(18 * time.Hour)  // 18:00

	sched, _ := schedule.NewSchedule(masterID, date, start, end, false)

	b1, _ := schedule.NewBreak(date.Add(13*time.Hour), date.Add(14*time.Hour)) // Break 13:00 - 14:00
	sched.AddBreak(b1)

	appt, _ := appointment.NewAppointment(uuid.New(), masterID, []uuid.UUID{}, date.Add(10*time.Hour), date.Add(11*time.Hour), 15*time.Minute) // 10:00 - 11:00 + 15m buffer

	appointments := []appointment.Appointment{*appt}

	windows := getFreeWindows(*sched, appointments)

	if len(windows) != 3 {
		t.Fatalf("expected 3 windows, got %d", len(windows))
	}

	// First window: 09:00 - 10:00
	if !windows[0].StartTime.Equal(start) || !windows[0].EndTime.Equal(date.Add(10*time.Hour)) {
		t.Errorf("window 0 wrong: %v - %v", windows[0].StartTime, windows[0].EndTime)
	}

	// Second window: 11:15 - 13:00
	if !windows[1].StartTime.Equal(date.Add(11*time.Hour+15*time.Minute)) || !windows[1].EndTime.Equal(date.Add(13*time.Hour)) {
		t.Errorf("window 1 wrong: %v - %v", windows[1].StartTime, windows[1].EndTime)
	}

	// Third window: 14:00 - 18:00
	if !windows[2].StartTime.Equal(date.Add(14*time.Hour)) || !windows[2].EndTime.Equal(end) {
		t.Errorf("window 2 wrong: %v - %v", windows[2].StartTime, windows[2].EndTime)
	}
}

func TestGroupIntoPeriods(t *testing.T) {
	date := time.Now().Truncate(24 * time.Hour)
	steps := []time.Time{
		date.Add(10 * time.Hour), // Morning
		date.Add(11 * time.Hour), // Morning
		date.Add(14 * time.Hour), // Day
		date.Add(18 * time.Hour), // Evening
	}

	periods := groupIntoPeriods(steps)

	if len(periods) != 3 {
		t.Fatalf("expected 3 periods, got %d", len(periods))
	}

	if periods[0].Name != "Утро" || len(periods[0].Slots) != 2 {
		t.Errorf("morning period wrong")
	}
	if periods[1].Name != "День" || len(periods[1].Slots) != 1 {
		t.Errorf("day period wrong")
	}
	if periods[2].Name != "Вечер" || len(periods[2].Slots) != 1 {
		t.Errorf("evening period wrong")
	}
}

func TestFilterWindows(t *testing.T) {
	now := time.Now()
	windows := []TimeSlot{
		{StartTime: now, EndTime: now.Add(30 * time.Minute)},
		{StartTime: now.Add(1 * time.Hour), EndTime: now.Add(2 * time.Hour)}, // 1 hour
	}

	filtered := filterWindows(windows, 45*time.Minute)

	if len(filtered) != 1 {
		t.Fatalf("expected 1 window, got %d", len(filtered))
	}

	if !filtered[0].StartTime.Equal(now.Add(1 * time.Hour)) {
		t.Errorf("wrong window filtered")
	}
}

func TestGenerateSteps(t *testing.T) {
	now := time.Now()
	windows := []TimeSlot{
		{StartTime: now, EndTime: now.Add(1 * time.Hour)},
	}

	steps := generateSteps(windows, 30*time.Minute, 15*time.Minute)

	// Should be: now, now+15m, now+30m
	if len(steps) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(steps))
	}

	expected := []time.Time{now, now.Add(15 * time.Minute), now.Add(30 * time.Minute)}
	for i, step := range steps {
		if !step.Equal(expected[i]) {
			t.Errorf("step %d wrong: got %v want %v", i, step, expected[i])
		}
	}
}
