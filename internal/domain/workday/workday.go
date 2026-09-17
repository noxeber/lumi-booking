package schedule

import (
	"errors"
	"time"
)

type ActionType string

const (
	ActionTypeBooking ActionType = "booking"
	ActionTypeBreak   ActionType = "break"
	ActionTypeLunch   ActionType = "lunch"
	ActionTypeDayOff  ActionType = "day_off"
)

var (
	ErrStartAfterEnd = errors.New("start time is after the end")
	ErrOverlap       = errors.New("new action overlaps with existing action")
	ErrOutWorkDay    = errors.New("action is out of work day bounds")
	ErrNilAction     = errors.New("action is nil")
)

type WorkDay struct {
	startTime time.Time
	endTime   time.Time
	actions   []Action
}

type Action struct {
	startTime  time.Time
	endTime    time.Time
	typeAction ActionType
}

type Interval struct {
	startTime time.Time
	endTime   time.Time
}

func NewWorkDay(startTime, endTime time.Time) (*WorkDay, error) {
	if !endTime.After(startTime) {
		return nil, ErrStartAfterEnd
	}
	return &WorkDay{startTime: startTime, endTime: endTime, actions: []Action{}}, nil
}

func NewAction(startTime, endTime time.Time, typeAction ActionType) (*Action, error) {
	if !endTime.After(startTime) {
		return nil, ErrStartAfterEnd
	}
	return &Action{startTime: startTime, endTime: endTime, typeAction: typeAction}, nil
}

func (w *WorkDay) AddAction(action *Action) error {
	if action == nil {
		return ErrNilAction
	}
	if action.startTime.Before(w.startTime) || action.endTime.After(w.endTime) {
		return ErrOutWorkDay
	}
	for _, existing := range w.actions {
		if action.startTime.Before(existing.endTime) &&
			existing.startTime.Before(action.endTime) {
			return ErrOverlap
		}
	}
	w.actions = append(w.actions, *action)
	return nil
}

func (w *WorkDay) Actions() []Action {
	return w.actions
}

func (w *WorkDay) FreeIntervals() []Interval {
	var free []Interval
	current := w.startTime
	for _, action := range w.actions {
		if action.startTime.After(current) {
			free = append(free, Interval{startTime: current, endTime: action.startTime})
		}
		if action.endTime.After(w.endTime) {
			free = append(free, Interval{startTime: current, endTime: w.endTime})
		}
		current = action.endTime
	}
	if current.Before(w.endTime) {
		free = append(free, Interval{startTime: current, endTime: w.endTime})
	}
	return free
}

func (a *Action) TypeAction() ActionType {
	return a.typeAction
}

func (i Interval) StartTime() time.Time { return i.startTime }
func (i Interval) EndTime() time.Time   { return i.endTime }

func (a Action) StartTime() time.Time { return a.startTime }
func (a Action) EndTime() time.Time   { return a.endTime }
