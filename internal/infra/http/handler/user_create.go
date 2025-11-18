package handler

import (
	"context"
	"net/http"

	user "scheduling/internal/app/user"
	infra "scheduling/internal/infra/gin"
)

type UserCreateHandler struct {
	UseCase *user.CreateUserUseCase
}

func NewUserCreateHandler(usecase *user.CreateUserUseCase) *UserCreateHandler {
	return &UserCreateHandler{UseCase: usecase}
}

func (handler *UserCreateHandler) Create(ctx infra.Context) error {

	name := ctx.Query("name");
	email := ctx.Query("email");
	password := ctx.Query("password")
	role := ctx.Query("role")

	ctxb := context.Background()

	input := user.UserInput{
		Name: name,
		Email: email,
		Password: password,
		Role: role,
	}
	users, err := handler.UseCase.Execute(ctxb, input)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError, 
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	return ctx.JSON(http.StatusOK, 
		map[string]interface{}{
			"usuario criado com success ": users,
		},
	)
}
