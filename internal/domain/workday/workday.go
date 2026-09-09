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
