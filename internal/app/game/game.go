package game

import (
	"database/sql"
	"fmt"
	sentryecho "github.com/getsentry/sentry-go/echo"
	_ "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/csrf"
	_csrfHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/csrf/delivery/http"
	_gameHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game/delivery/http"
	_gameRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game/repository"
	_gameUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game/usecase"
	_http "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/http"
	_ "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/logger"
	_middleware "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/middleware"
	_packHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack/delivery/http"
	_packRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack/repository"
	_sentry "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/sentry"
	_sessionRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/repository"
	_sessionUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/usecase"
	_userRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/microcosm-cc/bluemonday"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var log = logging.MustGetLogger("game")

func Start() {
	log.Info("Staring game service.")

	_sentry.InitSentry()

	e := echo.New()

	e.HTTPErrorHandler = _http.ErrorHandler

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

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(_middleware.NewMetricsMiddleware())

	e.Use(_middleware.LogMiddleware)

	//TODO: cleanup here
	e.Use(
		middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins: viper.GetStringSlice("server.CORS.allowed_hosts"),
				AllowHeaders: []string{
					echo.HeaderOrigin,
					echo.HeaderContentType,
					_csrfHttp.HeaderCSRFToken,
				},
				AllowCredentials: true,
			},
		),
	)

	csrfJWTToken := csrf.JwtToken{
		Secret: []byte(viper.GetString("server.CSRF.secret")),
	}

	dbDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.pass"),
		viper.GetString("database.db"),
	)
	db, err := sql.Open(
		"postgres",
		dbDSN,
	)
	if err != nil {
		log.Fatal("error init db: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("error connecting db: ", err)
	}
	defer db.Close()

	grpcConn, err := grpc.Dial(
		viper.GetString("server.session.host"),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("error dial to grpc service: ", err)
	}
	defer grpcConn.Close()

	sessionRepo := _sessionRepository.NewGrpcSessionRepository(grpcConn)
	sessionUseCase := _sessionUseCase.NewSessionUseCase(sessionRepo)

	userRepo := _userRepository.NewSqlUserRepository(db)

	packRepo := _packRepository.NewSqlPackRepository(db)
	packSanitizer := _packHttp.NewPackSanitizer(bluemonday.UGCPolicy())

	gameRepo := _gameRepository.NewMemGameRepository(userRepo)
	gameUseCase := _gameUseCase.NewGameUseCase(gameRepo, packRepo, packSanitizer)

	authMiddleware := _middleware.NewAuthMiddleware(sessionUseCase)

	csrfMiddleware := _middleware.NewCSRFMiddleware(csrfJWTToken)

	_csrfHttp.NewCSRFHandler(e, csrfJWTToken, authMiddleware)
	_gameHttp.NewGameHandler(e, gameUseCase, authMiddleware, csrfMiddleware)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
