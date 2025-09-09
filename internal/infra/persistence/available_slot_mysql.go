package persistence

import (
	"database/sql"
	"time"

	"scheduling/internal/domain/entities"
)

type AvailableSlotMySQLRepository struct {
	db *sql.DB
}

func NewAvailableSlotMySQLRepository(db *sql.DB) *AvailableSlotMySQLRepository {
	return &AvailableSlotMySQLRepository{db: db}
}

func (r *AvailableSlotMySQLRepository) FindByID(id int) (*entities.AvailableSlot, error) {
	query := "SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var slotID, staffID int
	var weekday string
	var startTime, endTime time.Time

	err := row.Scan(&slotID, &staffID, &weekday, &startTime, &endTime)
	if err != nil {
		return nil, err
	}

	slot, err := entities.NewAvailableSlot(staffID, entities.Weekday(weekday), startTime, endTime)
	if err != nil {
		return nil, err
	}
	slot.SetID(slotID)

	return slot, nil
}

func (r *AvailableSlotMySQLRepository) FindAllByStaffID(staffID int) ([]*entities.AvailableSlot, error) {
	query := "SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE staff_id = ?"
	rows, err := r.db.Query(query, staffID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*entities.AvailableSlot
	for rows.Next() {
		var slotID int
		var weekday string
		var startTime, endTime time.Time

		err := rows.Scan(&slotID, &staffID, &weekday, &startTime, &endTime)
		if err != nil {
			return nil, err
		}

		slot, err := entities.NewAvailableSlot(staffID, entities.Weekday(weekday), startTime, endTime)
		if err != nil {
			return nil, err
		}
		slot.SetID(slotID)
		slots = append(slots, slot)
	}

	return slots, nil
}

func (r *AvailableSlotMySQLRepository) FindSlotsByStaffAndDate(staffID int, date time.Time) ([]*entities.AvailableSlot, error) {
	weekday := entities.FromTimeWeekday(date.Weekday())
	return r.FindByWeekday(staffID, weekday)
}

func (r *AvailableSlotMySQLRepository) FindByWeekday(staffID int, weekday entities.Weekday) ([]*entities.AvailableSlot, error) {
	query := "SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE staff_id = ? AND weekday = ?"
	rows, err := r.db.Query(query, staffID, string(weekday))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*entities.AvailableSlot
	for rows.Next() {
		var slotID int
		var wd string
		var startTime, endTime time.Time

		err := rows.Scan(&slotID, &staffID, &wd, &startTime, &endTime)
		if err != nil {
			return nil, err
		}

		slot, err := entities.NewAvailableSlot(staffID, entities.Weekday(wd), startTime, endTime)
		if err != nil {
			return nil, err
		}
		slot.SetID(slotID)
		slots = append(slots, slot)
	}

	return slots, nil
}

func (r *AvailableSlotMySQLRepository) HasConflict(staffID int, weekday entities.Weekday, start, end time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) FROM available_slots
		WHERE staff_id = ? AND weekday = ? AND (
			(start_time < ? AND end_time > ?) OR
			(start_time >= ? AND start_time < ?)
		)
	`
	var count int
	err := r.db.QueryRow(query, staffID, string(weekday), end, start, start, end).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AvailableSlotMySQLRepository) IsWithinAvailableSlot(staffID int, start, end time.Time) (bool, error) {
	weekday := entities.FromTimeWeekday(start.Weekday())
	query := `
		SELECT COUNT(*) FROM available_slots
		WHERE staff_id = ? AND weekday = ? AND start_time <= ? AND end_time >= ?
	`
	var count int
	err := r.db.QueryRow(query, staffID, string(weekday), start, end).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AvailableSlotMySQLRepository) Save(slot *entities.AvailableSlot) error {
	query := "INSERT INTO available_slots (staff_id, weekday, start_time, end_time) VALUES (?, ?, ?, ?)"
	_, err := r.db.Exec(query,
		slot.StaffID(),
		string(slot.Weekday()),
		slot.StartTime(),
		slot.EndTime(),
	)
	return err
}

func (r *AvailableSlotMySQLRepository) Update(slot *entities.AvailableSlot) error {
	query := "UPDATE available_slots SET weekday = ?, start_time = ?, end_time = ? WHERE id = ?"
	_, err := r.db.Exec(query,
		string(slot.Weekday()),
		slot.StartTime(),
		slot.EndTime(),
		slot.ID(),
	)
	return err
}

func (r *AvailableSlotMySQLRepository) Delete(id int) error {
	query := "DELETE FROM available_slots WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
