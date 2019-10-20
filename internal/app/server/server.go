package server

import (
	_authHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/auth/delivery/http"
	_middleware "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/middleware"
	_sessionRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/repository"
	_sessionUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/usecase"
	_userHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/delivery/http"
	_userRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/repository"
	_userUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/usecase"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Start() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost", "http://frontend.photocouple.space"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType},
		AllowCredentials: true,
	}))

	userRepo := _userRepository.NewMemUserRepository()
	userUseCase := _userUseCase.NewUserUseCase(userRepo)

	sessionRepo := _sessionRepository.NewMemSessionRepository()
	sessionUseCase := _sessionUseCase.NewSessionUseCase(sessionRepo)

	authMiddleware := _middleware.NewAuthMiddleware(sessionUseCase)

	_userHttp.NewUserHandler(e, userUseCase, authMiddleware)
	_authHttp.NewAuthHandler(e, userUseCase, sessionUseCase, authMiddleware)

	e.Start("0.0.0.0:3000")
}
