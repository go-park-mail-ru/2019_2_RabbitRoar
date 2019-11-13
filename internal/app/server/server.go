package server

import (
	"context"
	"fmt"

	sentryecho "github.com/getsentry/sentry-go/echo"
	_authHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/auth/delivery/http"
	_ "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/csrf"
	_csrfHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/csrf/delivery/http"
	_gameHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game/delivery/http"
	_gameRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game/repository"
	_gameUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game/usecase"
	_ "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/logger"
	_middleware "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/middleware"
	_packHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack/delivery/http"
	_sentry "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/sentry"
	_sessionRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/repository"
	_sessionUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/usecase"
	_userHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/delivery/http"
	_userRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/repository"
	_userUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/usecase"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("server")

func Start() {
	log.Info("Staring service.")

	_sentry.InitSentry()

	e := echo.New()

	e.Use(_middleware.PanicMiddleware)

	e.Use(
		sentryecho.New(
			sentryecho.Options{
				Repanic:         true,
				WaitForDelivery: false,
				Timeout:         0,
			},
		),
	)

	e.Use(_middleware.LogMiddleware)

	e.Use(
		middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins: viper.GetStringSlice("server.CORS.allowed_hosts"),
				AllowHeaders: []string{
					echo.HeaderOrigin,
					echo.HeaderContentType,
					echo.HeaderUpgrade,
					"Connection",
					"Sec-WebSocket-Version",
					"Sec-WebSocket-Key",
					"Sec-WebSocket-Extensions",
					"Pragma",
					_csrfHttp.HeaderCSRFToken,
				},
				AllowCredentials: true,
			},
		),
	)

	jwtToken := csrf.JwtToken{
		Secret: []byte(viper.GetString("server.CSRF.secret")),
	}
	dbURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.pass"),
		viper.GetString("database.db"),
	)
	log.Info("dbURL: ", dbURL)
	pgxPool, err := pgxpool.Connect(
		context.Background(),
		dbURL,
	)
	if err != nil {
		log.Fatal("error connecting to db: ", err)
	}

	userRepo := _userRepository.NewSqlUserRepository(pgxPool)
	userUseCase := _userUseCase.NewUserUseCase(userRepo)

	sessionRepo := _sessionRepository.NewSqlSessionRepository(pgxPool)
	sessionUseCase := _sessionUseCase.NewSessionUseCase(sessionRepo)

	gameRepo := _gameRepository.NewSqlGameRepository(pgxPool)
	gameUseCase := _gameUseCase.NewGameUseCase(gameRepo)

	authMiddleware := _middleware.NewAuthMiddleware(sessionUseCase)
	csrfMiddleware := _middleware.NewCSRFMiddleware(jwtToken)

	_userHttp.NewUserHandler(e, userUseCase, authMiddleware, csrfMiddleware)
	_authHttp.NewAuthHandler(e, userUseCase, sessionUseCase, authMiddleware)
	_csrfHttp.NewCSRFHandler(e, jwtToken, authMiddleware)
	_packHttp.NewPackHandler(e, authMiddleware)
	_gameHttp.NewGameHandler(e, gameUseCase, authMiddleware, csrfMiddleware)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
