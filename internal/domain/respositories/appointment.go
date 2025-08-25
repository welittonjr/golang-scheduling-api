package respositories

import (
	"time"

	"scheduling/internal/domain/entities"
)

type AppointmentRepository interface {
	FindByID(id int) (*entities.Appointment, error)
	FindAllByStaffID(staffID int) ([]*entities.Appointment, error)
	HasConflict(staffID int, start, end time.Time) (bool, error)
	Save(appointment *entities.Appointment) error
	Update(appointment *entities.Appointment) error
	Delete(id int) error
}
