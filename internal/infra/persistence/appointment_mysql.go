package persistence

import (
	"database/sql"
	"time"

	"scheduling/internal/domain/entities"
)

type AppointmentMySQLRepository struct {
	db *sql.DB
}

func NewAppointmentMySQLRepository(db *sql.DB) *AppointmentMySQLRepository {
	return &AppointmentMySQLRepository{db: db}
}

func (r *AppointmentMySQLRepository) FindByID(id int) (*entities.Appointment, error) {
	query := "SELECT id, client_id, staff_id, service_id, scheduled_at, status, created_at FROM appointments WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var scheduledAt, createdAt time.Time
	var clientID, staffID, serviceID int
	var status string

	err := row.Scan(&id, &clientID, &staffID, &serviceID, &scheduledAt, &status, &createdAt)
	if err != nil {
		return nil, err
	}

	appointment, err := entities.NewAppointment(clientID, staffID, serviceID, scheduledAt)
	if err != nil {
		return nil, err
	}
	appointment.SetID(id)
	appointment.SetStatus(status)
	appointment.SetCreatedAt(createdAt)

	return appointment, nil
}

func (r *AppointmentMySQLRepository) FindAllByStaffID(staffID int) ([]*entities.Appointment, error) {
	query := "SELECT id, client_id, staff_id, service_id, scheduled_at, status, created_at FROM appointments WHERE staff_id = ?"
	rows, err := r.db.Query(query, staffID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*entities.Appointment
	for rows.Next() {
		var id, clientID, serviceID int
		var scheduledAt, createdAt time.Time
		var status string

		err := rows.Scan(&id, &clientID, &staffID, &serviceID, &scheduledAt, &status, &createdAt)
		if err != nil {
			return nil, err
		}

		appointment, err := entities.NewAppointment(clientID, staffID, serviceID, scheduledAt)
		if err != nil {
			return nil, err
		}
		appointment.SetID(id)
		appointment.SetStatus(status)
		appointment.SetCreatedAt(createdAt)

		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *AppointmentMySQLRepository) HasConflict(staffID int, start, end time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) FROM appointments
		WHERE staff_id = ? AND status = 'scheduled'
		AND (
			(scheduled_at BETWEEN ? AND ?)
			OR (? BETWEEN scheduled_at AND DATE_ADD(scheduled_at, INTERVAL 30 MINUTE))
		)
	`
	var count int
	err := r.db.QueryRow(query, staffID, start, end, end).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AppointmentMySQLRepository) Save(appointment *entities.Appointment) error {
	query := "INSERT INTO appointments (client_id, staff_id, service_id, scheduled_at, status, created_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query,
		appointment.ClientID(),
		appointment.StaffID(),
		appointment.ServiceID(),
		appointment.ScheduledAt(),
		appointment.Status(),
		appointment.CreatedAt(),
	)
	return err
}

func (r *AppointmentMySQLRepository) Update(appointment *entities.Appointment) error {
	query := "UPDATE appointments SET status = ? WHERE id = ?"
	_, err := r.db.Exec(query, appointment.Status(), appointment.ID())
	return err
}

func (r *AppointmentMySQLRepository) Delete(id int) error {
	query := "DELETE FROM appointments WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
