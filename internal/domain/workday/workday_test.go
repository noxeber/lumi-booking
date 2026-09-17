package schedule

import (
	"errors"
	"testing"
	"time"
)

// Helper functions
func mustNewWorkDay(t *testing.T, startTime, endTime time.Time) *WorkDay {
	t.Helper()
	wd, err := NewWorkDay(startTime, endTime)
	if err != nil {
		t.Fatalf("NewWorkDay() unexpected error: %v", err)
	}
	return wd
}

func mustNewAction(t *testing.T, startTime, endTime time.Time, actionType ActionType) *Action {
	t.Helper()
	action, err := NewAction(startTime, endTime, actionType)
	if err != nil {
		t.Fatalf("NewAction() unexpected error: %v", err)
	}
	return action
}

func timeMustParse(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02 15:04", value)
	if err != nil {
		t.Fatalf("time.Parse(%q) error: %v", value, err)
	}
	return parsed
}

func TestNewWorkDay(t *testing.T) {
	t.Run("creates workday with valid times", func(t *testing.T) {
		start := timeMustParse(t, "2026-09-18 09:00")
		end := timeMustParse(t, "2026-09-18 18:00")

		wd, err := NewWorkDay(start, end)
		if err != nil {
			t.Fatalf("NewWorkDay() unexpected error: %v", err)
		}

		if wd == nil {
			t.Fatal("NewWorkDay() returned nil")
		}

		if !wd.startTime.Equal(start) {
			t.Errorf("startTime = %v, want %v", wd.startTime, start)
		}

		if !wd.endTime.Equal(end) {
			t.Errorf("endTime = %v, want %v", wd.endTime, end)
		}

		if len(wd.actions) != 0 {
			t.Errorf("actions length = %d, want 0", len(wd.actions))
		}
	})

	t.Run("returns error when start equals end", func(t *testing.T) {
		time := timeMustParse(t, "2026-09-18 12:00")

		wd, err := NewWorkDay(time, time)

		if !errors.Is(err, ErrStartAfterEnd) {
			t.Errorf("NewWorkDay() error = %v, want %v", err, ErrStartAfterEnd)
		}

		if wd != nil {
			t.Error("NewWorkDay() should return nil on error")
		}
	})

	t.Run("returns error when start is after end", func(t *testing.T) {
		start := timeMustParse(t, "2026-09-18 18:00")
		end := timeMustParse(t, "2026-09-18 09:00")

		wd, err := NewWorkDay(start, end)

		if !errors.Is(err, ErrStartAfterEnd) {
			t.Errorf("NewWorkDay() error = %v, want %v", err, ErrStartAfterEnd)
		}

		if wd != nil {
			t.Error("NewWorkDay() should return nil on error")
		}
	})

	t.Run("creates workday with minimal duration", func(t *testing.T) {
		start := timeMustParse(t, "2026-09-18 09:00")
		end := start.Add(1 * time.Minute)

		wd, err := NewWorkDay(start, end)
		if err != nil {
			t.Fatalf("NewWorkDay() unexpected error: %v", err)
		}

		if wd == nil {
			t.Fatal("NewWorkDay() returned nil")
		}
	})
}

func TestNewAction(t *testing.T) {
	t.Run("creates action with valid times", func(t *testing.T) {
		start := timeMustParse(t, "2026-09-18 10:00")
		end := timeMustParse(t, "2026-09-18 11:00")
		actionType := ActionTypeBooking

		action, err := NewAction(start, end, actionType)
		if err != nil {
			t.Fatalf("NewAction() unexpected error: %v", err)
		}

		if action == nil {
			t.Fatal("NewAction() returned nil")
		}

		if !action.startTime.Equal(start) {
			t.Errorf("startTime = %v, want %v", action.startTime, start)
		}

		if !action.endTime.Equal(end) {
			t.Errorf("endTime = %v, want %v", action.endTime, end)
		}

		if action.typeAction != actionType {
			t.Errorf("typeAction = %v, want %v", action.typeAction, actionType)
		}
	})

	t.Run("creates action for all action types", func(t *testing.T) {
		start := timeMustParse(t, "2026-09-18 10:00")
		end := timeMustParse(t, "2026-09-18 11:00")

		testCases := []struct {
			actionType ActionType
			name       string
		}{
			{ActionTypeBooking, "booking"},
			{ActionTypeBreak, "break"},
			{ActionTypeLunch, "lunch"},
			{ActionTypeDayOff, "day_off"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				action, err := NewAction(start, end, tc.actionType)
				if err != nil {
					t.Fatalf("NewAction(%s) unexpected error: %v", tc.name, err)
				}

				if action.typeAction != tc.actionType {
					t.Errorf("typeAction = %v, want %v", action.typeAction, tc.actionType)
				}
			})
		}
	})

	t.Run("returns error when start equals end", func(t *testing.T) {
		time := timeMustParse(t, "2026-09-18 12:00")

		action, err := NewAction(time, time, ActionTypeBooking)

		if !errors.Is(err, ErrStartAfterEnd) {
			t.Errorf("NewAction() error = %v, want %v", err, ErrStartAfterEnd)
		}

		if action != nil {
			t.Error("NewAction() should return nil on error")
		}
	})

	t.Run("returns error when start is after end", func(t *testing.T) {
		start := timeMustParse(t, "2026-09-18 12:00")
		end := timeMustParse(t, "2026-09-18 11:00")

		action, err := NewAction(start, end, ActionTypeBooking)

		if !errors.Is(err, ErrStartAfterEnd) {
			t.Errorf("NewAction() error = %v, want %v", err, ErrStartAfterEnd)
		}

		if action != nil {
			t.Error("NewAction() should return nil on error")
		}
	})

	t.Run("creates action with minimal duration", func(t *testing.T) {
		start := timeMustParse(t, "2026-09-18 10:00")
		end := start.Add(1 * time.Second)

		action, err := NewAction(start, end, ActionTypeBreak)
		if err != nil {
			t.Fatalf("NewAction() unexpected error: %v", err)
		}

		if action == nil {
			t.Fatal("NewAction() returned nil")
		}
	})
}

func TestAddAction(t *testing.T) {
	t.Run("adds valid action successfully", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		actionStart := timeMustParse(t, "2026-09-18 10:00")
		actionEnd := timeMustParse(t, "2026-09-18 11:00")
		action := mustNewAction(t, actionStart, actionEnd, ActionTypeBooking)

		err := wd.AddAction(action)
		if err != nil {
			t.Fatalf("AddAction() unexpected error: %v", err)
		}

		if len(wd.actions) != 1 {
			t.Fatalf("actions length = %d, want 1", len(wd.actions))
		}

		if !wd.actions[0].startTime.Equal(actionStart) {
			t.Errorf("actions[0].startTime = %v, want %v", wd.actions[0].startTime, actionStart)
		}
	})

	t.Run("returns error for nil action", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		err := wd.AddAction(nil)

		if !errors.Is(err, ErrNilAction) {
			t.Errorf("AddAction(nil) error = %v, want %v", err, ErrNilAction)
		}

		if len(wd.actions) != 0 {
			t.Errorf("actions length = %d, want 0", len(wd.actions))
		}
	})

	t.Run("returns error when action starts before workday", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		actionStart := timeMustParse(t, "2026-09-18 08:00")
		actionEnd := timeMustParse(t, "2026-09-18 10:00")
		action := mustNewAction(t, actionStart, actionEnd, ActionTypeBooking)

		err := wd.AddAction(action)

		if !errors.Is(err, ErrOutWorkDay) {
			t.Errorf("AddAction() error = %v, want %v", err, ErrOutWorkDay)
		}

		if len(wd.actions) != 0 {
			t.Errorf("actions length = %d, want 0", len(wd.actions))
		}
	})

	t.Run("returns error when action ends after workday", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		actionStart := timeMustParse(t, "2026-09-18 17:00")
		actionEnd := timeMustParse(t, "2026-09-18 19:00")
		action := mustNewAction(t, actionStart, actionEnd, ActionTypeBooking)

		err := wd.AddAction(action)

		if !errors.Is(err, ErrOutWorkDay) {
			t.Errorf("AddAction() error = %v, want %v", err, ErrOutWorkDay)
		}

		if len(wd.actions) != 0 {
			t.Errorf("actions length = %d, want 0", len(wd.actions))
		}
	})

	t.Run("returns error when action completely outside workday", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		actionStart := timeMustParse(t, "2026-09-18 19:00")
		actionEnd := timeMustParse(t, "2026-09-18 20:00")
		action := mustNewAction(t, actionStart, actionEnd, ActionTypeBooking)

		err := wd.AddAction(action)

		if !errors.Is(err, ErrOutWorkDay) {
			t.Errorf("AddAction() error = %v, want %v", err, ErrOutWorkDay)
		}
	})

	t.Run("accepts action at workday boundaries", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		// Action starts exactly at workday start
		action1Start := wdStart
		action1End := timeMustParse(t, "2026-09-18 10:00")
		action1 := mustNewAction(t, action1Start, action1End, ActionTypeBooking)

		err := wd.AddAction(action1)
		if err != nil {
			t.Fatalf("AddAction() at start boundary unexpected error: %v", err)
		}

		// Action ends exactly at workday end
		action2Start := timeMustParse(t, "2026-09-18 17:00")
		action2End := wdEnd
		action2 := mustNewAction(t, action2Start, action2End, ActionTypeBreak)

		err = wd.AddAction(action2)
		if err != nil {
			t.Fatalf("AddAction() at end boundary unexpected error: %v", err)
		}

		if len(wd.actions) != 2 {
			t.Errorf("actions length = %d, want 2", len(wd.actions))
		}
	})

	t.Run("returns error when action overlaps with existing", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		// Add first action: 10:00 - 11:00
		action1Start := timeMustParse(t, "2026-09-18 10:00")
		action1End := timeMustParse(t, "2026-09-18 11:00")
		action1 := mustNewAction(t, action1Start, action1End, ActionTypeBooking)
		wd.AddAction(action1)

		// Try to add overlapping action: 10:30 - 11:30
		action2Start := timeMustParse(t, "2026-09-18 10:30")
		action2End := timeMustParse(t, "2026-09-18 11:30")
		action2 := mustNewAction(t, action2Start, action2End, ActionTypeBreak)

		err := wd.AddAction(action2)

		if !errors.Is(err, ErrOverlap) {
			t.Errorf("AddAction() error = %v, want %v", err, ErrOverlap)
		}

		if len(wd.actions) != 1 {
			t.Errorf("actions length = %d, want 1", len(wd.actions))
		}
	})

	t.Run("returns error when action completely contains existing", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		// Add first action: 11:00 - 12:00
		action1Start := timeMustParse(t, "2026-09-18 11:00")
		action1End := timeMustParse(t, "2026-09-18 12:00")
		action1 := mustNewAction(t, action1Start, action1End, ActionTypeBooking)
		wd.AddAction(action1)

		// Try to add containing action: 10:00 - 13:00
		action2Start := timeMustParse(t, "2026-09-18 10:00")
		action2End := timeMustParse(t, "2026-09-18 13:00")
		action2 := mustNewAction(t, action2Start, action2End, ActionTypeBreak)

		err := wd.AddAction(action2)

		if !errors.Is(err, ErrOverlap) {
			t.Errorf("AddAction() error = %v, want %v", err, ErrOverlap)
		}
	})

	t.Run("accepts action that starts exactly when previous ends", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		// First action: 10:00 - 11:00
		action1Start := timeMustParse(t, "2026-09-18 10:00")
		action1End := timeMustParse(t, "2026-09-18 11:00")
		action1 := mustNewAction(t, action1Start, action1End, ActionTypeBooking)
		wd.AddAction(action1)

		// Second action starts exactly at 11:00
		action2Start := timeMustParse(t, "2026-09-18 11:00")
		action2End := timeMustParse(t, "2026-09-18 12:00")
		action2 := mustNewAction(t, action2Start, action2End, ActionTypeBreak)

		err := wd.AddAction(action2)
		if err != nil {
			t.Fatalf("AddAction() unexpected error: %v", err)
		}

		if len(wd.actions) != 2 {
			t.Errorf("actions length = %d, want 2", len(wd.actions))
		}
	})

	t.Run("adds multiple non-overlapping actions", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		actions := []struct {
			start string
			end   string
		}{
			{"2026-09-18 10:00", "2026-09-18 11:00"},
			{"2026-09-18 12:00", "2026-09-18 13:00"},
			{"2026-09-18 14:00", "2026-09-18 15:00"},
		}

		for _, a := range actions {
			start := timeMustParse(t, a.start)
			end := timeMustParse(t, a.end)
			action := mustNewAction(t, start, end, ActionTypeBooking)

			err := wd.AddAction(action)
			if err != nil {
				t.Fatalf("AddAction(%s-%s) unexpected error: %v", a.start, a.end, err)
			}
		}

		if len(wd.actions) != 3 {
			t.Errorf("actions length = %d, want 3", len(wd.actions))
		}
	})

	t.Run("detects overlap with any existing action", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		// Add three actions
		wd.AddAction(mustNewAction(t, timeMustParse(t, "2026-09-18 10:00"), timeMustParse(t, "2026-09-18 11:00"), ActionTypeBooking))
		wd.AddAction(mustNewAction(t, timeMustParse(t, "2026-09-18 12:00"), timeMustParse(t, "2026-09-18 13:00"), ActionTypeBreak))
		wd.AddAction(mustNewAction(t, timeMustParse(t, "2026-09-18 14:00"), timeMustParse(t, "2026-09-18 15:00"), ActionTypeLunch))

		// Try to overlap with middle action
		overlappingAction := mustNewAction(t, timeMustParse(t, "2026-09-18 12:30"), timeMustParse(t, "2026-09-18 13:30"), ActionTypeBooking)

		err := wd.AddAction(overlappingAction)

		if !errors.Is(err, ErrOverlap) {
			t.Errorf("AddAction() error = %v, want %v", err, ErrOverlap)
		}
	})
}

func TestActions(t *testing.T) {
	t.Run("returns empty slice for new workday", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		actions := wd.Actions()

		if len(actions) != 0 {
			t.Errorf("Actions() length = %d, want 0", len(actions))
		}
	})

	t.Run("returns all added actions", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		action1 := mustNewAction(t, timeMustParse(t, "2026-09-18 10:00"), timeMustParse(t, "2026-09-18 11:00"), ActionTypeBooking)
		action2 := mustNewAction(t, timeMustParse(t, "2026-09-18 12:00"), timeMustParse(t, "2026-09-18 13:00"), ActionTypeBreak)

		wd.AddAction(action1)
		wd.AddAction(action2)

		actions := wd.Actions()

		if len(actions) != 2 {
			t.Fatalf("Actions() length = %d, want 2", len(actions))
		}
	})

	t.Run("returns copy of actions slice", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		action := mustNewAction(t, timeMustParse(t, "2026-09-18 10:00"), timeMustParse(t, "2026-09-18 11:00"), ActionTypeBooking)
		wd.AddAction(action)

		actions := wd.Actions()
		actions[0] = Action{} // Modify returned slice

		// Original should be unchanged
		originalActions := wd.Actions()
		if len(originalActions) != 1 {
			t.Error("modifying returned slice should not affect original")
		}
	})
}

func TestFreeIntervals(t *testing.T) {
	t.Run("returns one interval for empty workday", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		intervals := wd.FreeIntervals()

		if len(intervals) != 1 {
			t.Fatalf("FreeIntervals() length = %d, want 1", len(intervals))
		}

		if !intervals[0].startTime.Equal(wdStart) {
			t.Errorf("intervals[0].startTime = %v, want %v", intervals[0].startTime, wdStart)
		}

		if !intervals[0].endTime.Equal(wdEnd) {
			t.Errorf("intervals[0].endTime = %v, want %v", intervals[0].endTime, wdEnd)
		}
	})

	t.Run("returns two intervals when action in middle", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		actionStart := timeMustParse(t, "2026-09-18 12:00")
		actionEnd := timeMustParse(t, "2026-09-18 13:00")
		wd.AddAction(mustNewAction(t, actionStart, actionEnd, ActionTypeBooking))

		intervals := wd.FreeIntervals()

		if len(intervals) != 2 {
			t.Fatalf("FreeIntervals() length = %d, want 2", len(intervals))
		}

		// First interval: 09:00 - 12:00
		if !intervals[0].startTime.Equal(wdStart) {
			t.Errorf("intervals[0].startTime = %v, want %v", intervals[0].startTime, wdStart)
		}
		if !intervals[0].endTime.Equal(actionStart) {
			t.Errorf("intervals[0].endTime = %v, want %v", intervals[0].endTime, actionStart)
		}

		// Second interval: 13:00 - 18:00
		if !intervals[1].startTime.Equal(actionEnd) {
			t.Errorf("intervals[1].startTime = %v, want %v", intervals[1].startTime, actionEnd)
		}
		if !intervals[1].endTime.Equal(wdEnd) {
			t.Errorf("intervals[1].endTime = %v, want %v", intervals[1].endTime, wdEnd)
		}
	})

	t.Run("returns one interval when action at start", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		actionStart := wdStart
		actionEnd := timeMustParse(t, "2026-09-18 10:00")
		wd.AddAction(mustNewAction(t, actionStart, actionEnd, ActionTypeBreak))

		intervals := wd.FreeIntervals()

		if len(intervals) != 1 {
			t.Fatalf("FreeIntervals() length = %d, want 1", len(intervals))
		}

		if !intervals[0].startTime.Equal(actionEnd) {
			t.Errorf("intervals[0].startTime = %v, want %v", intervals[0].startTime, actionEnd)
		}
		if !intervals[0].endTime.Equal(wdEnd) {
			t.Errorf("intervals[0].endTime = %v, want %v", intervals[0].endTime, wdEnd)
		}
	})

	t.Run("returns one interval when action at end", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		actionStart := timeMustParse(t, "2026-09-18 17:00")
		actionEnd := wdEnd
		wd.AddAction(mustNewAction(t, actionStart, actionEnd, ActionTypeLunch))

		intervals := wd.FreeIntervals()

		if len(intervals) != 1 {
			t.Fatalf("FreeIntervals() length = %d, want 1", len(intervals))
		}

		if !intervals[0].startTime.Equal(wdStart) {
			t.Errorf("intervals[0].startTime = %v, want %v", intervals[0].startTime, wdStart)
		}
		if !intervals[0].endTime.Equal(actionStart) {
			t.Errorf("intervals[0].endTime = %v, want %v", intervals[0].endTime, actionStart)
		}
	})

	t.Run("returns multiple intervals for multiple actions", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		// Add three actions with gaps
		wd.AddAction(mustNewAction(t, timeMustParse(t, "2026-09-18 10:00"), timeMustParse(t, "2026-09-18 11:00"), ActionTypeBooking))
		wd.AddAction(mustNewAction(t, timeMustParse(t, "2026-09-18 12:00"), timeMustParse(t, "2026-09-18 13:00"), ActionTypeBreak))
		wd.AddAction(mustNewAction(t, timeMustParse(t, "2026-09-18 14:00"), timeMustParse(t, "2026-09-18 15:00"), ActionTypeLunch))

		intervals := wd.FreeIntervals()

		// Should have 4 free intervals
		if len(intervals) != 4 {
			t.Fatalf("FreeIntervals() length = %d, want 4", len(intervals))
		}

		expected := []struct {
			start string
			end   string
		}{
			{"2026-09-18 09:00", "2026-09-18 10:00"},
			{"2026-09-18 11:00", "2026-09-18 12:00"},
			{"2026-09-18 13:00", "2026-09-18 14:00"},
			{"2026-09-18 15:00", "2026-09-18 18:00"},
		}

		for i, exp := range expected {
			start := timeMustParse(t, exp.start)
			end := timeMustParse(t, exp.end)

			if !intervals[i].startTime.Equal(start) {
				t.Errorf("intervals[%d].startTime = %v, want %v", i, intervals[i].startTime, start)
			}
			if !intervals[i].endTime.Equal(end) {
				t.Errorf("intervals[%d].endTime = %v, want %v", i, intervals[i].endTime, end)
			}
		}
	})

	t.Run("returns no intervals when actions are consecutive", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		// Add consecutive actions covering entire workday
		wd.AddAction(mustNewAction(t, wdStart, timeMustParse(t, "2026-09-18 12:00"), ActionTypeBooking))
		wd.AddAction(mustNewAction(t, timeMustParse(t, "2026-09-18 12:00"), wdEnd, ActionTypeBreak))

		intervals := wd.FreeIntervals()

		if len(intervals) != 0 {
			t.Errorf("FreeIntervals() length = %d, want 0", len(intervals))
		}
	})

	t.Run("returns empty slice when entire workday is covered", func(t *testing.T) {
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		// Single action covering entire workday
		wd.AddAction(mustNewAction(t, wdStart, wdEnd, ActionTypeDayOff))

		intervals := wd.FreeIntervals()

		if len(intervals) != 0 {
			t.Errorf("FreeIntervals() length = %d, want 0", len(intervals))
		}
	})
}

func TestActionTypeAction(t *testing.T) {
	t.Run("returns correct action type", func(t *testing.T) {
		start := timeMustParse(t, "2026-09-18 10:00")
		end := timeMustParse(t, "2026-09-18 11:00")

		testCases := []struct {
			actionType ActionType
			name       string
		}{
			{ActionTypeBooking, "booking"},
			{ActionTypeBreak, "break"},
			{ActionTypeLunch, "lunch"},
			{ActionTypeDayOff, "day_off"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				action, _ := NewAction(start, end, tc.actionType)

				got := action.TypeAction()

				if got != tc.actionType {
					t.Errorf("TypeAction() = %v, want %v", got, tc.actionType)
				}
			})
		}
	})
}

func TestIntervalGetters(t *testing.T) {
	t.Run("returns correct start and end times", func(t *testing.T) {
		start := timeMustParse(t, "2026-09-18 10:00")
		end := timeMustParse(t, "2026-09-18 11:00")

		interval := Interval{startTime: start, endTime: end}

		if !interval.StartTime().Equal(start) {
			t.Errorf("StartTime() = %v, want %v", interval.StartTime(), start)
		}

		if !interval.EndTime().Equal(end) {
			t.Errorf("EndTime() = %v, want %v", interval.EndTime(), end)
		}
	})

	t.Run("returns zero times for zero interval", func(t *testing.T) {
		var interval Interval

		if !interval.StartTime().IsZero() {
			t.Errorf("StartTime() = %v, want zero time", interval.StartTime())
		}

		if !interval.EndTime().IsZero() {
			t.Errorf("EndTime() = %v, want zero time", interval.EndTime())
		}
	})
}

func TestActionGetters(t *testing.T) {
	t.Run("returns correct start and end times", func(t *testing.T) {
		start := timeMustParse(t, "2026-09-18 10:00")
		end := timeMustParse(t, "2026-09-18 11:00")

		action, _ := NewAction(start, end, ActionTypeBooking)

		if !action.StartTime().Equal(start) {
			t.Errorf("StartTime() = %v, want %v", action.StartTime(), start)
		}

		if !action.EndTime().Equal(end) {
			t.Errorf("EndTime() = %v, want %v", action.EndTime(), end)
		}
	})

	t.Run("returns zero times for zero action", func(t *testing.T) {
		var action Action

		if !action.StartTime().IsZero() {
			t.Errorf("StartTime() = %v, want zero time", action.StartTime())
		}

		if !action.EndTime().IsZero() {
			t.Errorf("EndTime() = %v, want zero time", action.EndTime())
		}
	})
}

func TestErrorMessages(t *testing.T) {
	t.Run("ErrStartAfterEnd has correct message", func(t *testing.T) {
		expected := "start time is after the end"
		if ErrStartAfterEnd.Error() != expected {
			t.Errorf("ErrStartAfterEnd.Error() = %q, want %q", ErrStartAfterEnd.Error(), expected)
		}
	})

	t.Run("ErrOverlap has correct message", func(t *testing.T) {
		expected := "new action overlaps with existing action"
		if ErrOverlap.Error() != expected {
			t.Errorf("ErrOverlap.Error() = %q, want %q", ErrOverlap.Error(), expected)
		}
	})

	t.Run("ErrOutWorkDay has correct message", func(t *testing.T) {
		expected := "action is out of work day bounds"
		if ErrOutWorkDay.Error() != expected {
			t.Errorf("ErrOutWorkDay.Error() = %q, want %q", ErrOutWorkDay.Error(), expected)
		}
	})

	t.Run("ErrNilAction has correct message", func(t *testing.T) {
		expected := "action is nil"
		if ErrNilAction.Error() != expected {
			t.Errorf("ErrNilAction.Error() = %q, want %q", ErrNilAction.Error(), expected)
		}
	})
}

func TestScheduleWorkflow(t *testing.T) {
	t.Run("complete workday workflow", func(t *testing.T) {
		// Create workday
		wdStart := timeMustParse(t, "2026-09-18 09:00")
		wdEnd := timeMustParse(t, "2026-09-18 18:00")
		wd := mustNewWorkDay(t, wdStart, wdEnd)

		// Add lunch break
		lunchStart := timeMustParse(t, "2026-09-18 13:00")
		lunchEnd := timeMustParse(t, "2026-09-18 14:00")
		lunch := mustNewAction(t, lunchStart, lunchEnd, ActionTypeLunch)

		err := wd.AddAction(lunch)
		if err != nil {
			t.Fatalf("AddAction(lunch) unexpected error: %v", err)
		}

		// Add booking before lunch
		booking1Start := timeMustParse(t, "2026-09-18 10:00")
		booking1End := timeMustParse(t, "2026-09-18 11:30")
		booking1 := mustNewAction(t, booking1Start, booking1End, ActionTypeBooking)

		err = wd.AddAction(booking1)
		if err != nil {
			t.Fatalf("AddAction(booking1) unexpected error: %v", err)
		}

		// Add booking after lunch
		booking2Start := timeMustParse(t, "2026-09-18 14:30")
		booking2End := timeMustParse(t, "2026-09-18 16:00")
		booking2 := mustNewAction(t, booking2Start, booking2End, ActionTypeBooking)

		err = wd.AddAction(booking2)
		if err != nil {
			t.Fatalf("AddAction(booking2) unexpected error: %v", err)
		}

		// Try to add overlapping booking
		overlappingStart := timeMustParse(t, "2026-09-18 11:00")
		overlappingEnd := timeMustParse(t, "2026-09-18 12:00")
		overlapping := mustNewAction(t, overlappingStart, overlappingEnd, ActionTypeBooking)

		err = wd.AddAction(overlapping)
		if !errors.Is(err, ErrOverlap) {
			t.Errorf("AddAction(overlapping) error = %v, want %v", err, ErrOverlap)
		}

		// Check actions count
		actions := wd.Actions()
		if len(actions) != 3 {
			t.Errorf("Actions() length = %d, want 3", len(actions))
		}

		// Check free intervals
		intervals := wd.FreeIntervals()
		// Expected: 09:00-10:00, 11:30-13:00, 14:00-14:30, 16:00-18:00
		if len(intervals) != 4 {
			t.Errorf("FreeIntervals() length = %d, want 4", len(intervals))
		}
	})
}
