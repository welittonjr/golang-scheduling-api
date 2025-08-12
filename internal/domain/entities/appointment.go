package entities

import (
	"errors"
	"time"
)

type AppointmentStatus string

const (
	StatusScheduled AppointmentStatus = "scheduled"
	StatusCompleted AppointmentStatus = "completed"
	StatusCancelled AppointmentStatus = "cancelled"
)

type Appointment struct {
	id          int
	clientID    int
	staffID     int
	serviceID   int
	scheduledAt time.Time
	status      AppointmentStatus
	createdAt   time.Time
}

func NewAppointment(clientID, staffID, serviceID int, scheduledAt time.Time) (*Appointment, error) {
	if clientID == 0 || staffID == 0 || serviceID == 0 {
		return nil, errors.New("cliente, profissional e serviço são obrigatórios")
	}
	if scheduledAt.Before(time.Now()) {
		return nil, errors.New("não é possível agendar para o passado")
	}

	return &Appointment{
		clientID:    clientID,
		staffID:     staffID,
		serviceID:   serviceID,
		scheduledAt: scheduledAt,
		status:      StatusScheduled,
		createdAt:   time.Now(),
	}, nil
}

func (a *Appointment) Cancel() {
	a.status = StatusCancelled
}

func (a *Appointment) Complete() {
	a.status = StatusCompleted
}

func (a *Appointment) IsScheduled() bool {
	return a.status == StatusScheduled
}

func (a *Appointment) IsCanceled() bool {
	return a.status == StatusCancelled
}

func (a *Appointment) SetID(id int) {
	a.id = id
}

func (a *Appointment) SetStatus(status string) {
	switch AppointmentStatus(status) {
	case StatusScheduled, StatusCompleted, StatusCancelled:
		a.status = AppointmentStatus(status)
	default:
		a.status = StatusScheduled
	}
}

func (a *Appointment) SetCreatedAt(t time.Time) {
	a.createdAt = t
}

func (a *Appointment) ID() int                   { return a.id }
func (a *Appointment) ClientID() int             { return a.clientID }
func (a *Appointment) StaffID() int              { return a.staffID }
func (a *Appointment) ServiceID() int            { return a.serviceID }
func (a *Appointment) ScheduledAt() time.Time    { return a.scheduledAt }
func (a *Appointment) Status() AppointmentStatus { return a.status }
func (a *Appointment) CreatedAt() time.Time      { return a.createdAt }
