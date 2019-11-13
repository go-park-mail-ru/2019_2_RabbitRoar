package http

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

type handler struct {}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewGameHandler(
	e *echo.Echo,
	authMiddleware echo.MiddlewareFunc,
) {
	handler := handler{}

	group := e.Group("/game")
	group.GET("/ws", handler.ws)
}

func (h *handler) ws(ctx echo.Context) error {
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			log.Error(err)
		}

		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Error(err)
		}

		err = ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Error(err)
		}
	}
}
