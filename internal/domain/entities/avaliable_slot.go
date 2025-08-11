package entities

import (
	"errors"
	"fmt"
	"time"
)

type Weekday string

const (
	Sunday    Weekday = "sunday"
	Monday    Weekday = "monday"
	Tuesday   Weekday = "tuesday"
	Wednesday Weekday = "wednesday"
	Thursday  Weekday = "thursday"
	Friday    Weekday = "friday"
	Saturday  Weekday = "saturday"
)

type AvailableSlot struct {
	id        int
	staffID   int
	weekday   Weekday
	startTime time.Time
	endTime   time.Time
}

func NewAvailableSlot(staffID int, weekday Weekday, start, end time.Time) (*AvailableSlot, error) {
	if staffID == 0 {
		return nil, errors.New("staffID é obrigatório")
	}
	if !isValidWeekday(weekday) {
		return nil, fmt.Errorf("dia da semana inválido: %s", weekday)
	}
	if !start.Before(end) {
		return nil, errors.New("o horário inicial deve ser antes do final")
	}

	return &AvailableSlot{
		staffID:   staffID,
		weekday:   weekday,
		startTime: start,
		endTime:   end,
	}, nil
}

func isValidWeekday(day Weekday) bool {
	switch day {
	case Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday:
		return true
	default:
		return false
	}
}

func FromTimeWeekday(w time.Weekday) Weekday {
	switch w {
	case time.Sunday:
		return Sunday
	case time.Monday:
		return Monday
	case time.Tuesday:
		return Tuesday
	case time.Wednesday:
		return Wednesday
	case time.Thursday:
		return Thursday
	case time.Friday:
		return Friday
	case time.Saturday:
		return Saturday
	default:
		return Weekday("")
	}
}

func (s *AvailableSlot) SetID(id int)         { s.id = id }
func (s *AvailableSlot) ID() int              { return s.id }
func (s *AvailableSlot) StaffID() int         { return s.staffID }
func (s *AvailableSlot) Weekday() Weekday     { return s.weekday }
func (s *AvailableSlot) StartTime() time.Time { return s.startTime }
func (s *AvailableSlot) EndTime() time.Time   { return s.endTime }
