package http

import (
	"net/http"
	"sync"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/http_utils"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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

	group := e.Group("/game", authMiddleware)
	group.GET("/", handler.self)
	group.POST("/", csrfMiddleware(handler.create))
	group.POST("/:uuid/join", handler.join)
	group.DELETE("/leave", handler.leave)
	group.GET("/ws", handler.ws)
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

	creator := ctx.Get("user").(*models.User)

	if err := gh.usecase.SQLCreate(g, *creator); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := gh.usecase.MemCreate(g, *creator); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
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

	if err := gh.usecase.JoinPlayerToGame(userID, gameID); err != nil {
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
		return err
	}

	defer ws.Close()

	user := ctx.Get("user").(*models.User)

	conn := gh.usecase.NewConnection()

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

	var wg sync.WaitGroup
	wg.Add(2)

	go conn.RunReceive(ws, &wg)
	go conn.RunSend(ws, &wg)

	wg.Wait()

	return nil
}
