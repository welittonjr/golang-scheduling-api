package services

import (
	"context"
	"fmt"
	"log/slog"
	"scheduling/internal/domain/entities"
	"scheduling/internal/domain/repositories"
	"time"
)

type UserService struct {
	logger   *slog.Logger
	userRepo repositories.UserRepository
}

func NewUserService(
	logger *slog.Logger,
	userRepo repositories.UserRepository,
) *UserService {
	return &UserService{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (userService *UserService) Create(ctx context.Context, user entities.User) error {

	startTime := time.Now()
	
	exist, err := userService.userRepo.EmailExist(ctx, user.Email())
	if err != nil {
		userService.logger.Error(
			"Erro ao verificar existencia do email do usuário",
			"error", err.Error(),
			"email", user.Email(),
			"operation", "user_service.check_email",
			"duration_ms", time.Since(startTime).Milliseconds(),
		)
		return err
	}
	
	if exist {
		userService.logger.Warn("tentativa de criar usuário com email existente",
			"email",  user.Email(),
			"operation", "user_service.duplicate_email",
			"duration_ms", time.Since(startTime).Milliseconds(),
		)
		return fmt.Errorf("email já existe: %s", user.Email())
	}

	err = userService.userRepo.Create(ctx, &user)
	if err != nil {
		userService.logger.Error(
			"Erro ao tentar criar o usuário",
			"payload",user,
			"operation", "user_service.create_user",
			"durations_ms", time.Since(startTime).Milliseconds(),
		)
	}

	return nil
}

func (userService *UserService) Authentication(ctx context.Context, user entities.User) (bool, error) {
	
	userDB, err := userService.userRepo.FindByEmail(ctx, user.Email())	
	if err != nil {
		userService.logger.Error(
			"Erro ao buscar usuário por email",
			"error", err.Error(),
			"email", user.Email(),
			"operation", "user_service.fidByEmail",
		)
		return false, err
	}

	if !user.CheckPassword(userDB.Password()) {
		return false, nil
	}
	
	return true, nil
}