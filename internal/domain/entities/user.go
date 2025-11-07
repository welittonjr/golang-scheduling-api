package entities

import (
	"errors"
	"scheduling/internal/domain/valueobject"
	"time"
)

const (
	RoleClient = "client"
	RoleAdmin  = "admin"
)

type User struct {
	id        int
	name      string
	email     valueobject.Email
	password  string
	role      string
	createdAt time.Time
}

func NewUser(id int, name string, emailStr, password, role string) (*User, error) {
	if name == "" {
		return nil, errors.New("nome é obrigatório")
	}

	if len(password) < 6 {
		return nil, errors.New("a senha deve ter ao menos 6 caracteres")
	}
	if role != RoleClient && role != RoleAdmin {
		return nil, errors.New("papel inválido")
	}

	email, err := valueobject.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	return &User{
		id:        id,
		name:      name,
		email:     email,
		password:  password,
		role:      role,
		createdAt: time.Now(),
	}, nil
}

func RebuildUser(id int, name string, email valueobject.Email, role string) *User {
	return &User{
		id:    id,
		name:  name,
		email: email,
		role:  role,
	}
}

func (u *User) SetCreatedAt(t time.Time) {
	u.createdAt = t
}

func (u *User) CanAccessAdminPanel() bool {
	return u.role == RoleAdmin
}

func (u *User) IsClient() bool {
	return u.role == RoleClient
}

func (u *User) CheckPassword(password string) bool {
	return u.password == password
}

func (u *User) ID() int              { return u.id }
func (u *User) Name() string         { return u.name }
func (u *User) Email() string        { return u.email.String() }
func (u *User) Password() string     { return u.password }
func (u *User) Role() string         { return u.role }
func (u *User) CreatedAt() time.Time { return u.createdAt }
