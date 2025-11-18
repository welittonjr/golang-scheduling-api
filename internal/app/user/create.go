package user

import (
	"context"
	"scheduling/internal/domain/entities"
	"scheduling/internal/domain/services"
)

type CreateUserUseCase struct {
	UserService *services.UserService
}

func NewCreateUserUseCase(
	userSerivce *services.UserService,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserService: userSerivce,
	}
}

func (useCase *CreateUserUseCase) Execute(ctx context.Context, input UserInput) (*UserOutput, error) {
	
	user, err := entities.NewUser(
		input.ID,
		input.Name,
		input.Email,
		input.Password,
		input.Role,
	)

	if err != nil {
		return nil, err
	}

	err = useCase.UserService.Create(ctx, *user)
	if err != nil {
		return  nil, err
	}

	return &UserOutput{
		ID:       user.ID(),
		Name:     user.Name(),
	}, nil
}