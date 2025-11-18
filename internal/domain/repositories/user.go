package repositories

import (
	"context"
	"scheduling/internal/domain/entities"
)

type UserRepository interface {
	Repository[entities.User]
	Exists(ctx context.Context, id int) (bool, error)
	EmailExist(ctx context.Context, email string) (bool, error)
	FindByEmail(ctx context.Context, email string) (entities.User, error)
}
