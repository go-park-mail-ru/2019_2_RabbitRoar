package http

import (
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
	group.GET("", handler.self)
	group.POST("", authMiddleware(csrfMiddleware(handler.create)))
	group.POST("/:uuid/join", authMiddleware(csrfMiddleware(handler.join)))
	group.DELETE("/leave", authMiddleware(csrfMiddleware(handler.leave)))
	group.GET("/ws", authMiddleware(handler.ws))
}

func (gh *handler) self(ctx echo.Context) error {
	page := http_utils.GetIntParam(ctx, 0)

	if page < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "page less than 0 provided")
	}

	content, err := gh.usecase.Fetch(page)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "error fetching page of games",
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
			Message:  "can't parse game object",
			Internal: err,
		}
	}

	if g.PlayersCapacity > viper.GetInt("internal.players_cap_limit") {
		return echo.NewHTTPError(http.StatusBadRequest, "players capacity is too big")
	}

	creator := ctx.Get("user").(*models.User)

	if err := gh.usecase.Create(g, *creator); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.NoContent(http.StatusCreated)
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

	userID := ctx.Get("user").(*models.User).ID

	g, err := gh.usecase.GetByID(gameID)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "error finding the game",
			Internal: err,
		}
	}

	if g, err = gh.usecase.JoinPlayerToGame(userID, gameID); err != nil {
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
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusUpgradeRequired,
			Message:  "error establishing websocket connection",
			Internal: err,
		}
	}

	user := ctx.Get("user").(*models.User)

	conn := gh.usecase.NewConnectionWrapper(ws)

	gameID, err := gh.usecase.GetGameIDByUserID(user.ID)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "error finding game ID",
			Internal: err,
		}
	}

	err = gh.usecase.JoinConnectionToGame(gameID, *user, conn)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "unable to join the game",
			Internal: err,
		}
	}

	go conn.RunReceive(user.ID)
	go conn.RunSend()

	return ctx.NoContent(http.StatusOK)
}
