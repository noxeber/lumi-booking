package schedule

import (
	"errors"
	"time"
	"uuid"
)

var (
	ErrStartAfterEnd = errors.New("start time is after the end")
	ErrOverlap       = errors.New("new break overlaps with existing break")
	ErrOutSchedule   = errors.New("break is out of schedule bounds")
)

type Break struct {
	startTime time.Time
	endTime   time.Time
}

type Schedule struct {
	masterID  uuid.UUID
	date      time.Time
	startTime time.Time
	endTime   time.Time
	breaks    []Break
	isDayOff  bool
}

func NewSchedule(masterID uuid.UUID, date time.Time, startTime, endTime time.Time, isDayOff bool) (*Schedule, error) {
	if !isDayOff && !endTime.After(startTime) {
		return nil, ErrStartAfterEnd
	}
	return &Schedule{masterID: masterID, date: date, startTime: startTime, endTime: endTime, isDayOff: isDayOff}, nil
}

func NewBreak(startTime, endTime time.Time) (*Break, error) {
	if !endTime.After(startTime) {
		return nil, ErrStartAfterEnd
	}
	return &Break{startTime: startTime, endTime: endTime}, nil
}

func (s *Schedule) AddBreak(b *Break) error {
	if s.isDayOff {
		return nil // Cannot add break to a day off, or we could return an error, but doing nothing is fine
	}
	if b.startTime.Before(s.startTime) || b.endTime.After(s.endTime) {
		return ErrOutSchedule
	}
	for _, existing := range s.breaks {
		if b.startTime.Before(existing.endTime) && existing.startTime.Before(b.endTime) {
			return ErrOverlap
		}
	}
	s.breaks = append(s.breaks, *b)
	return nil
}

func (s *Schedule) MasterID() uuid.UUID  { return s.masterID }
func (s *Schedule) Date() time.Time      { return s.date }
func (s *Schedule) StartTime() time.Time { return s.startTime }
func (s *Schedule) EndTime() time.Time   { return s.endTime }
func (s *Schedule) Breaks() []Break      { return s.breaks }
func (s *Schedule) IsDayOff() bool       { return s.isDayOff }

func (b Break) StartTime() time.Time { return b.startTime }
func (b Break) EndTime() time.Time   { return b.endTime }
