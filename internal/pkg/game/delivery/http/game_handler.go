package http

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack"
	"github.com/prometheus/common/log"
	"net/http"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/http_utils"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type handler struct {
	usecase game.UseCase
	packSanitizer pack.Sanitizer
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewGameHandler(
	e *echo.Echo,
	uc game.UseCase,
	authMiddleware echo.MiddlewareFunc,
	csrfMiddleware echo.MiddlewareFunc,
) {
	handler := handler{
		usecase: uc,
	}

	group := e.Group("/game")
	group.GET("", handler.list)
	group.POST("", authMiddleware(csrfMiddleware(handler.create)))
	group.POST("/:uuid/join", authMiddleware(csrfMiddleware(handler.join)))
	group.DELETE("/leave", authMiddleware(csrfMiddleware(handler.leave)))
	group.GET("/ws", authMiddleware(handler.ws))
}

func (gh *handler) list(ctx echo.Context) error {
	page := http_utils.GetIntParam(ctx, "page", 0)

	if page < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "page less than 0 provided")
	}

	content, err := gh.usecase.Fetch(page)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "unable to fetch a page of games",
			Internal: err,
		}
	}

	return ctx.JSON(http.StatusOK, *content)
}

func (gh *handler) create(ctx echo.Context) error {
	var g models.Game
	err := ctx.Bind(&g)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "unable to identify game object",
			Internal: err,
		}
	}

	if g.PlayersCapacity > viper.GetInt("internal.players_cap_limit") {
		return echo.NewHTTPError(http.StatusBadRequest, "players capacity is too big")
	}

	creator := ctx.Get("user").(*models.User)

	if err := gh.usecase.Create(&g, *creator); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "unable to create a game",
			Internal: err,
		}
	}

	return ctx.JSON(http.StatusCreated, g)
}

func (gh *handler) join(ctx echo.Context) error {
	gameID, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUnprocessableEntity,
			Message:  "can't parse game uuid",
			Internal: err,
		}
	}

	u := ctx.Get("user").(*models.User)

	var g *models.Game
	if g, err = gh.usecase.JoinPlayerToGame(*u, gameID); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "error joining the game",
			Internal: err,
		}
	}

	return ctx.JSON(http.StatusOK, *g)
}

func (gh *handler) leave(ctx echo.Context) error {
	userID := ctx.Get("user").(*models.User).ID

	if err := gh.usecase.KickPlayerFromGame(userID); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "error leaving the game",
			Internal: err,
		}
	}

	return ctx.NoContent(http.StatusOK)
}

func (gh *handler) ws(ctx echo.Context) error {
	log.Info("ws header Host: ", ctx.Request().Header.Get("Host"))
	log.Info("ws header Origin: ", ctx.Request().Header.Get("Origin"))
	log.Info("ws header Connection: ", ctx.Request().Header.Get("Connection"))
	log.Info("ws header Upgrade: ", ctx.Request().Header.Get("Upgrade"))
	log.Info("ws header Sec-WebSocket-Key: ", ctx.Request().Header.Get("Sec-WebSocket-Key"))
	log.Info("ws header Sec-WebSocket-Version: ", ctx.Request().Header.Get("Sec-WebSocket-Version"))

	userID := ctx.Get("user").(*models.User).ID

	gameID, err := gh.usecase.GetGameIDByUserID(userID)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "error finding game ID",
			Internal: err,
		}
	}

	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUpgradeRequired,
			Message:  "error establishing websocket connection",
			Internal: err,
		}
	}

	conn := gh.usecase.NewConnectionWrapper(ws)

	err = gh.usecase.JoinConnectionToGame(gameID, userID, conn)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "unable to join the game",
			Internal: err,
		}
	}

	go conn.RunReceive(userID)
	go conn.RunSend()
	log.Info("WS handler reached the end.")
	return nil
}
