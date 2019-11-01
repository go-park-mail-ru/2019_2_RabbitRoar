package server

import (
	_authHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/auth/delivery/http"
	_ "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/csrf"
	_ "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/logger"
	_middleware "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/middleware"
	_sessionRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/repository"
	_sessionUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/usecase"
	_userHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/delivery/http"
	_userRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/repository"
	_userUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/usecase"
	_csrfHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/csrf/delivery/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("server")

func Start() {
	log.Info("Staring service.")
	e := echo.New()

	e.Use(_middleware.PanicMiddleware)
	e.Use(
		middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins:     viper.GetStringSlice("server.CORS.allowed_hosts"),
				AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType},
				AllowCredentials: true,
			},
		),
	)

	jwtToken := csrf.JwtToken{
		Secret: []byte(viper.GetString("server.CSRF.secret")),
	}

	userRepo := _userRepository.NewMemUserRepository()
	userUseCase := _userUseCase.NewUserUseCase(userRepo)

	sessionRepo := _sessionRepository.NewMemSessionRepository()
	sessionUseCase := _sessionUseCase.NewSessionUseCase(sessionRepo)

	authMiddleware := _middleware.NewAuthMiddleware(sessionUseCase)
	csrfMiddleware := _middleware.NewCSRFMiddleware(jwtToken)

	_userHttp.NewUserHandler(e, userUseCase, authMiddleware, csrfMiddleware)
	_authHttp.NewAuthHandler(e, userUseCase, sessionUseCase, authMiddleware)
	_csrfHttp.NewCSRFHandler(e, jwtToken, authMiddleware)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
