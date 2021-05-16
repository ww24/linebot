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
	// Next returns next scheduled time or ErrEndSchedule.
	Next(t time.Time) (time.Time, error)
}

func ParseScheduler(serialized string) (Scheduler, error) {
	c := strings.SplitN(serialized, schedulerSep, 2)
	switch schedulerType(c[0]) {
	case schedulerTypeOnetime:
		s := &OnetimeScheduler{}
		if err := s.parse(serialized); err != nil {
			return nil, err
		}
		return s, nil
	case schedulerTypeDaily:
		s := &DailyScheduler{}
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
	schedulerTypeOnetime schedulerType = "o"
	schedulerTypeDaily   schedulerType = "d"

	schedulerSep = "#"
)

type OnetimeScheduler struct {
	Time time.Time
}

func (s *OnetimeScheduler) parse(serialized string) error {
	st := strings.TrimPrefix(serialized, schedulerTypeOnetime.String()+schedulerSep)
	t, err := time.Parse(time.RFC3339, st)
	if err != nil {
		return xerrors.Errorf("failed to parse time: %w", err)
	}

	s.Time = t
	return nil
}

func (s *OnetimeScheduler) String() string {
	return schedulerTypeOnetime.String() + schedulerSep + s.Time.Format(time.RFC3339)
}

func (s *OnetimeScheduler) Next(t time.Time) (time.Time, error) {
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

func (s *DailyScheduler) Next(t time.Time) (time.Time, error) {
	year, month, day := t.Date()
	hour, min, sec := s.Time.Clock()
	target := time.Date(year, month, day, hour, min, sec, 0, s.Time.Location())
	if t.Before(target) {
		return target, nil
	}
	return target.AddDate(0, 0, 1), nil
}
