package repositories

import "scheduling/internal/domain/entities"

type UserRepository interface {
	FindByID(id int) (*entities.User, error)
	Exists(id int) (bool, error)
}
