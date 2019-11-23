package chat

import (
	"database/sql"
	"fmt"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/chat"
	_chatHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/chat/delivery/http"
	_ "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/config"
	_csrfHttp "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/csrf/delivery/http"
	_http "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/http"
	_ "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/logger"
	_middleware "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/middleware"
	_sessionRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/repository"
	_sessionUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("chat server")

func Start() {
	log.Info("Staring chat service.")

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

	//e.Use(_middleware.NewMetricsMiddleware())

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

	sessionRepo := _sessionRepository.NewSqlSessionRepository(db)
	sessionUseCase := _sessionUseCase.NewSessionUseCase(sessionRepo)

	authMiddleware := _middleware.NewAuthMiddleware(sessionUseCase)

	_chatHttp.NewChatHandler(e, sessionUseCase, authMiddleware, chat.NewHub())

	log.Fatal(e.Start(viper.GetString("chat.address")))
}
