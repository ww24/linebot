package model

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/xerrors"
)

var (
	ErrEndSchedule          = errors.New("end of schedule")
	ErrInvalidSchedulerType = errors.New("invalid scheduler type")
)

type Scheduler interface {
	// String returns a serialized schedule settings.
	String() string
	// UIText returns a text for UI.
	UIText() string
	// Next returns next scheduled time or ErrEndSchedule.
	Next(t time.Time) (time.Time, error)
}

func ParseScheduler(serialized string) (Scheduler, error) {
	c := strings.SplitN(serialized, schedulerSep, schedulerElementSize)
	switch schedulerType(c[0]) {
	case schedulerTypeOneshot:
		s := new(OneshotScheduler)
		if err := s.parse(serialized); err != nil {
			return nil, err
		}
		return s, nil
	case schedulerTypeDaily:
		s := new(DailyScheduler)
		if err := s.parse(serialized); err != nil {
			return nil, err
		}
		return s, nil
	}
	return nil, ErrInvalidSchedulerType
}

type schedulerType string

func (t schedulerType) String() string {
	return string(t)
}

const (
	schedulerTypeOneshot schedulerType = "o"
	schedulerTypeDaily   schedulerType = "d"

	schedulerSep         = "#"
	schedulerElementSize = 2
)

type OneshotScheduler struct {
	Time time.Time
}

func (s *OneshotScheduler) parse(serialized string) error {
	st := strings.TrimPrefix(serialized, schedulerTypeOneshot.String()+schedulerSep)
	t, err := time.Parse(time.RFC3339, st)
	if err != nil {
		return xerrors.Errorf("failed to parse time: %w", err)
	}

	s.Time = t
	return nil
}

func (s *OneshotScheduler) String() string {
	return schedulerTypeOneshot.String() + schedulerSep + s.Time.Format(time.RFC3339)
}

func (s *OneshotScheduler) UIText() string {
	return s.Time.Format("at 2006-01-02 15:04.")
}

func (s *OneshotScheduler) Next(t time.Time) (time.Time, error) {
	if !t.Before(s.Time) {
		return time.Time{}, ErrEndSchedule
	}
	return s.Time, nil
}

type DailyScheduler struct {
	Time time.Time
}

func (s *DailyScheduler) parse(serialized string) error {
	st := strings.TrimPrefix(serialized, schedulerTypeDaily.String()+schedulerSep)
	t, err := time.Parse(time.RFC3339, st)
	if err != nil {
		return xerrors.Errorf("failed to parse time: %w", err)
	}

	s.Time = t
	return nil
}

func (s *DailyScheduler) String() string {
	return schedulerTypeDaily.String() + schedulerSep + s.Time.Format(time.RFC3339)
}

func (s *DailyScheduler) UIText() string {
	return s.Time.Format("at 15:04 every day.")
}

func (s *DailyScheduler) Next(t time.Time) (time.Time, error) {
	loc := s.Time.Location()
	year, month, day := t.In(loc).Date()
	hour, min, sec := s.Time.Clock()
	target := time.Date(year, month, day, hour, min, sec, 0, loc)
	if t.Before(target) {
		return target, nil
	}
	return target.AddDate(0, 0, 1), nil
}
