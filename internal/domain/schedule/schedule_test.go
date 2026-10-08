package schedule

import (
	"errors"
	"testing"
	"time"

	"uuid"
)

func TestNewSchedule(t *testing.T) {
	masterID := uuid.New()
	date := time.Now()

	t.Run("creates valid schedule", func(t *testing.T) {
		start := time.Now()
		end := start.Add(8 * time.Hour)

		s, err := NewSchedule(masterID, date, start, end, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if s.StartTime() != start {
			t.Errorf("StartTime = %v, want %v", s.StartTime(), start)
		}
		if s.EndTime() != end {
			t.Errorf("EndTime = %v, want %v", s.EndTime(), end)
		}
		if s.IsDayOff() {
			t.Errorf("IsDayOff = true, want false")
		}
	})

	t.Run("fails when end before start on workday", func(t *testing.T) {
		start := time.Now()
		end := start.Add(-1 * time.Hour)

		_, err := NewSchedule(masterID, date, start, end, false)
		if !errors.Is(err, ErrStartAfterEnd) {
			t.Errorf("error = %v, want %v", err, ErrStartAfterEnd)
		}
	})

	t.Run("allows end before start on day off", func(t *testing.T) {
		start := time.Now()
		end := start.Add(-1 * time.Hour)

		s, err := NewSchedule(masterID, date, start, end, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !s.IsDayOff() {
			t.Errorf("IsDayOff = false, want true")
		}
	})
}

func TestNewBreak(t *testing.T) {
	t.Run("creates valid break", func(t *testing.T) {
		start := time.Now()
		end := start.Add(1 * time.Hour)

		b, err := NewBreak(start, end)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if b.StartTime() != start || b.EndTime() != end {
			t.Errorf("break mismatch")
		}
	})

	t.Run("fails when end before start", func(t *testing.T) {
		start := time.Now()
		end := start.Add(-1 * time.Hour)

		_, err := NewBreak(start, end)
		if !errors.Is(err, ErrStartAfterEnd) {
			t.Errorf("error = %v, want %v", err, ErrStartAfterEnd)
		}
	})
}

func TestAddBreak(t *testing.T) {
	masterID := uuid.New()
	date := time.Now()
	start := time.Now()
	end := start.Add(8 * time.Hour)

	t.Run("adds break within schedule", func(t *testing.T) {
		s, _ := NewSchedule(masterID, date, start, end, false)
		b, _ := NewBreak(start.Add(4*time.Hour), start.Add(5*time.Hour))

		err := s.AddBreak(b)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(s.Breaks()) != 1 {
			t.Errorf("len = %d, want 1", len(s.Breaks()))
		}
	})

	t.Run("fails when break out of bounds", func(t *testing.T) {
		s, _ := NewSchedule(masterID, date, start, end, false)
		b, _ := NewBreak(start.Add(9*time.Hour), start.Add(10*time.Hour))

		err := s.AddBreak(b)
		if !errors.Is(err, ErrOutSchedule) {
			t.Errorf("error = %v, want %v", err, ErrOutSchedule)
		}
	})

	t.Run("fails when breaks overlap", func(t *testing.T) {
		s, _ := NewSchedule(masterID, date, start, end, false)
		b1, _ := NewBreak(start.Add(4*time.Hour), start.Add(5*time.Hour))
		b2, _ := NewBreak(start.Add(4*time.Hour+30*time.Minute), start.Add(5*time.Hour+30*time.Minute))

		s.AddBreak(b1)
		err := s.AddBreak(b2)
		if !errors.Is(err, ErrOverlap) {
			t.Errorf("error = %v, want %v", err, ErrOverlap)
		}
	})
}
