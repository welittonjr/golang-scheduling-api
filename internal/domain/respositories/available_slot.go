package respositories

import (
	"time"

	"scheduling/internal/domain/entities"
)

type AvailableSlotRepository interface {
	FindByID(id int) (*entities.AvailableSlot, error)
	FindAllByStaffID(staffID int) ([]*entities.AvailableSlot, error)
	FindSlotsByStaffAndDate(staffID int, date time.Time) ([]*entities.AvailableSlot, error)
	HasConflict(staffID int, weekday entities.Weekday, start, end time.Time) (bool, error)
	IsWithinAvailableSlot(staffID int, start, end time.Time) (bool, error)
	Save(slot *entities.AvailableSlot) error
	Update(slot *entities.AvailableSlot) error
	Delete(id int) error
}
