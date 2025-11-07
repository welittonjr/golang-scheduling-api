package persistence

import (
	"context"
	"database/sql"
	"scheduling/internal/domain/entities"
	"scheduling/internal/domain/valueobject"
	"time"
)

type UserMySQLRepository struct {
	db *sql.DB
}

func NewUserMySQLRepository(db *sql.DB) *UserMySQLRepository {
	return &UserMySQLRepository{db: db}
}

func (r *UserMySQLRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (name, email, password, role, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Name(),
		user.Email(),
		user.Password(),
		user.Role(),
		now,
	)

	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (r *UserMySQLRepository) FindByID(ctx context.Context, id int) (*entities.User, error) {
	query := `
		SELECT id, name, email, password, role, created_at 
		FROM users 
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var userID int
	var name, email, password, role string
	var createdAt sql.NullTime

	err := row.Scan(
		&userID, 
		&name, 
		&email, 
		&password, 
		&role, 
		&createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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

func (r *UserMySQLRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET name = ?, email = ?, role = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Name(),
		user.Email(),
		user.Role(),
		user.ID(),
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *UserMySQLRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE users 
		SET deleted_at = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, now, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *UserMySQLRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT id, name, email, role, created_at, updated_at 
		FROM users 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var userID int
		var name, email, role string
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(&userID, &name, &email, &role, &createdAt, &updatedAt)
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

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserMySQLRepository) Count(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL"
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserMySQLRepository) Exists(ctx context.Context, id int) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE id = ?"
	var count int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
