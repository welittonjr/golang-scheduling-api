package main

import (
	_ "github.com/go-sql-driver/mysql"

	"scheduling/internal/infra/database"
	http "scheduling/internal/infra/gin"
	ginadapter "scheduling/internal/infra/gin/adapter"
	"scheduling/internal/infra/logger"
	"scheduling/internal/infra/middleware"

	"scheduling/internal/app/user"
	"scheduling/internal/domain/services"
	"scheduling/internal/infra/http/handler"
	"scheduling/internal/infra/persistence"
)

func main() {

	logger := logger.SetupLogger();

	db := database.Connect();

	userRepo := persistence.NewUserMySQLRepository(db)

	userService := services.NewUserService(logger, userRepo)

	userUseCase := user.NewCreateUserUseCase(userService)

	userHandler := handler.NewUserCreateHandler(userUseCase)

	router := ginadapter.NewRouter()

	router.Use(middleware.TraceIDMiddleware())

	router.GET("/", func(ctx http.Context) error {
		return ctx.JSON(200, map[string]string{"message": "Alive S2!"})
	})

	router.POST("/user", userHandler.Create)

	router.Run(":8080")
}
