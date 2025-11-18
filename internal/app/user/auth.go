package user

import (
	"context"
	"fmt"
	"scheduling/internal/domain/entities"
	"scheduling/internal/domain/services"
	"scheduling/internal/infra/jwt"
)


type AuthUseCase struct {
	UserService *services.UserService
}

func NewAuthUseCase(
	userSerivce *services.UserService,
) *AuthUseCase {
	return &AuthUseCase{
		UserService: userSerivce,
	}
}

func (useCase AuthUseCase) Execute(ctx context.Context, input UserAuthInput) (*UserAuthOutput, error) {

	user, err := entities.NewUser(
		0,
		"",
		input.Email,
		input.Password,
		"",
	)

	if err != nil {
		return nil, err
	}

	ok, err := useCase.UserService.Authentication(ctx, *user)
	if err != nil {
		return  nil, err
	}

	if !ok {
		return &UserAuthOutput{
			Token: "",
		}, nil
	}

	token, err := jwt.CreateToken(user.Name())
	if err != nil {
		return  nil, fmt.Errorf("Erro ao criar o token de autenticação %s", err.Error())
	}

	return &UserAuthOutput{
		Token: token,
	}, nil
}
