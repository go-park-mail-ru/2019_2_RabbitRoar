package server

import (
	"database/sql"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"io/ioutil"

	sentryecho "github.com/getsentry/sentry-go/echo"
	_authHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/auth/delivery/http"
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
	_packUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack/usecase"
	_sentry "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/sentry"
	_sessionRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/repository"
	_sessionUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/usecase"
	_userHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/delivery/http"
	_userRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/repository"
	_userUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
)

var log = logging.MustGetLogger("server")

func Start() {
	log.Info("Staring service.")

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
					echo.HeaderUpgrade,
					"Connection",
					"Sec-WebSocket-Version",
					"Sec-WebSocket-Key",
					"Sec-WebSocket-Extensions",
					_csrfHttp.HeaderCSRFToken,
				},
				AllowCredentials: true,
			},
		),
	)

	jwtToken := csrf.JwtToken{
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
	log.Info("dbURL: ", dbDSN)
	db, err := sql.Open(
		"postgres",
		dbDSN,
	)
	if err != nil {
		log.Fatal("error connecting to db: ", err)
	}

	userRepo := _userRepository.NewSqlUserRepository(db)
	userUseCase := _userUseCase.NewUserUseCase(userRepo)

	sessionRepo := _sessionRepository.NewSqlSessionRepository(db)
	sessionUseCase := _sessionUseCase.NewSessionUseCase(sessionRepo)

	schemaBytes, err := ioutil.ReadFile(viper.GetString("server.schema.pack"))
	if err != nil {
		log.Fatal("error reading schema for pack", err)
	}
	packSchema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(schemaBytes))
	if err != nil {
		log.Fatal("error parsing schema for pack", err)
	}
	packRepo := _packRepository.NewSqlPackRepository(db)
	packUseCase := _packUseCase.NewUserUseCase(packRepo)
	packSanitizer := _packHttp.NewPackSanitizer(bluemonday.UGCPolicy())

	gameMemRepo := _gameRepository.NewMemGameRepository()
	gameUseCase := _gameUseCase.NewGameUseCase(gameMemRepo, packRepo, packSanitizer)

	authMiddleware := _middleware.NewAuthMiddleware(sessionUseCase)

	csrfMiddleware := _middleware.NewCSRFMiddleware(jwtToken)

	_userHttp.NewUserHandler(e, userUseCase, authMiddleware, csrfMiddleware)
	_authHttp.NewAuthHandler(e, userUseCase, sessionUseCase, authMiddleware)
	_csrfHttp.NewCSRFHandler(e, jwtToken, authMiddleware)
	_gameHttp.NewGameHandler(e, gameUseCase, authMiddleware, csrfMiddleware)
	_packHttp.NewPackHandler(e, packUseCase, userUseCase,  sessionUseCase, authMiddleware, csrfMiddleware, packSchema)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
