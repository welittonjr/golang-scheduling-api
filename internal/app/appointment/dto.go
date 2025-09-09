package appointment

import (
	"time"
	"scheduling/internal/domain/entities"
)

type AppointmentInput struct {
	AppointmentID int    `json:"appointment_id,omitempty"`
	ClientID      int    `json:"client_id"`
	StaffID       int    `json:"staff_id"`
	ServiceID     int    `json:"service_id"`
	ScheduledAt   string `json:"scheduled_at"`
}

type AppointmentOutput struct {
	ID          int       `json:"id"`
	ClientID    int       `json:"client_id"`
	StaffID     int       `json:"staff_id"`
	ServiceID   int       `json:"service_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewAppointmentOutput(appointment *entities.Appointment) *AppointmentOutput {
	return &AppointmentOutput{
		ID:          appointment.ID(),
		ClientID:    appointment.ClientID(),
		StaffID:     appointment.StaffID(),
		ServiceID:   appointment.ServiceID(),
		ScheduledAt: appointment.ScheduledAt(),
		Status:      string(appointment.Status()),
		CreatedAt:   appointment.CreatedAt(),
	}
}