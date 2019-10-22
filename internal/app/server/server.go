package server

import (
	"fmt"
	_authHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/auth/delivery/http"
	_middleware "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/middleware"
	_sessionRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/repository"
	_sessionUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/usecase"
	_userHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/delivery/http"
	_userRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/repository"
	_userUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/usecase"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.SetConfigType("json")
	viper.SetConfigName("server")
	viper.AddConfigPath("configs")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	fmt.Println("Config file used: ", viper.ConfigFileUsed())
}

func Start() {
	e := echo.New()

	e.Use(_middleware.PanicMiddleware)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     viper.GetStringSlice("server.CORS.allowed_hosts"),
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

	log.Fatal(e.Start(viper.GetString("server.address")))
}
