package persistence

import (
	"database/sql"

	"scheduling/internal/domain/entities"
	"scheduling/internal/domain/valueobject"
)

type UserMySQLRepository struct {
	db *sql.DB
}

func NewUserMySQLRepository(db *sql.DB) *UserMySQLRepository {
	return &UserMySQLRepository{db: db}
}

func (r *UserMySQLRepository) FindByID(id int) (*entities.User, error) {
	query := "SELECT id, name, email, role, created_at FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var userID int
	var name, email, role string
	var createdAt sql.NullTime

	err := row.Scan(&userID, &name, &email, &role, &createdAt)
	if err != nil {
		return nil, err
	}

	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	user := entities.RebuildUser(userID, name, emailVO, role)

	if createdAt.Valid {
		user.SetCreatedAt(createdAt.Time)
	}

	return user, nil
}

func (r *UserMySQLRepository) Exists(id int) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE id = ?"
	var count int
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
