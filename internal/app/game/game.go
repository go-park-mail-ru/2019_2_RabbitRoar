package game

import (
	"database/sql"
	"encoding/json"
	"flag"
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
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/sd"
	_sentry "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/sentry"
	_sessionRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/repository"
	_sessionUseCase "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/microcosm-cc/bluemonday"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/balancer"
	"strconv"
)

var log = logging.MustGetLogger("game")

func Start() {
	gameLimiter := sd.GameLimiter{}
	if err := gameLimiter.Init(); err != nil {
		log.Fatal("error consul error: ", err)
	}

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

	// prikol begin
	flag.Parse()

	config := consulapi.DefaultConfig()
	config.Address = *consulAddr
	consul, err = consulapi.NewClient(config)

	health, _, err := consul.Health().Service("session-api", "", false, nil)
	if err != nil {
		log.Fatalf("cant get alive services")
	}

	servers := []string{}
	for _, item := range health {
		addr := item.Service.Address +
			":" + strconv.Itoa(item.Service.Port)
		servers = append(servers, addr)
	}

	nameResolver = &testNameResolver{
		addr: servers[0],
	}
	//

	grpcConn, err := grpc.Dial(
		servers[0],
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithBalancer(grpc.RoundRobin(nameResolver)),
	)

	if err != nil {
		log.Fatal("error dial to grpc service: ", err)
	}
	defer grpcConn.Close()

	if len(servers) > 1 {
		var updates []*naming.Update
		for i := 1; i < len(servers); i++ {
			updates = append(updates, &naming.Update{
				Op:   naming.Add,
				Addr: servers[i],
			})
		}
		nameResolver.w.inject(updates)
	}

	go runOnlineServiceDiscovery(servers)

	sessionRepo := _sessionRepository.NewGrpcSessionRepository(grpcConn)
	sessionUseCase := _sessionUseCase.NewSessionUseCase(sessionRepo)

	packRepo := _packRepository.NewSqlPackRepository(db)
	packSanitizer := _packHttp.NewPackSanitizer(bluemonday.UGCPolicy())

	gameMemRepo := _gameRepository.NewMemGameRepository()
	gameUseCase := _gameUseCase.NewGameUseCase(gameMemRepo, packRepo, packSanitizer, &gameLimiter)

	authMiddleware := _middleware.NewAuthMiddleware(sessionUseCase)

	csrfMiddleware := _middleware.NewCSRFMiddleware(csrfJWTToken)

	_csrfHttp.NewCSRFHandler(e, csrfJWTToken, authMiddleware)
	_gameHttp.NewGameHandler(e, gameUseCase, authMiddleware, csrfMiddleware)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
