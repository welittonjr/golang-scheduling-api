package persistence

import (
	"database/sql"

	"scheduling/internal/domain/entities"
)

type ServiceMySQLRepository struct {
	db *sql.DB
}

func NewServiceMySQLRepository(db *sql.DB) *ServiceMySQLRepository {
	return &ServiceMySQLRepository{db: db}
}

func (r *ServiceMySQLRepository) FindByID(id int) (*entities.Service, error) {
	query := "SELECT id, staff_id, name, duration, price FROM services WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var serviceID, staffID, duration int
	var name string
	var price float64

	err := row.Scan(&serviceID, &staffID, &name, &duration, &price)
	if err != nil {
		return nil, err
	}

	service, err := entities.NewService(serviceID, staffID, name, duration, price)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (r *ServiceMySQLRepository) FindAllByStaffID(staffID int) ([]*entities.Service, error) {
	query := "SELECT id, staff_id, name, duration, price FROM services WHERE staff_id = ?"
	rows, err := r.db.Query(query, staffID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*entities.Service
	for rows.Next() {
		var id, duration int
		var name string
		var price float64

		err := rows.Scan(&id, &staffID, &name, &duration, &price)
		if err != nil {
			return nil, err
		}

		service, err := entities.NewService(id, staffID, name, duration, price)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, nil
}

func (r *ServiceMySQLRepository) Exists(id int) (bool, error) {
	query := "SELECT COUNT(*) FROM services WHERE id = ?"
	var count int
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
