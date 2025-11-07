package repositories

import (
	"scheduling/internal/domain/entities"
)

type UserRepository interface {
	Repository[entities.User]
	Exists(id int) (bool, error)
}
