package server

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/delivery/http"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/repository"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/usecase"
	"github.com/labstack/echo"
)

func Start() {
	e := echo.New()

	userRepo := repository.NewMemUserRepository()
	userUseCase := usecase.NewUserUseCase(userRepo)
	userGroup := e.Group("/user")
	http.NewUserHandler(userGroup, userUseCase)

	e.Start("0.0.0.0:3000")
}
