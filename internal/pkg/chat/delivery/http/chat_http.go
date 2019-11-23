package http

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/chat"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session"
	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
	"log"
)

type handler struct {
	sessionUseCase session.UseCase
	sanitizer      *bluemonday.Policy
	hub *chat.Hub
}

func NewChatHandler(
	e *echo.Echo,
	sessionUseCase session.UseCase,
	authMiddleware echo.MiddlewareFunc,
	hub *chat.Hub,
) {
	handler := handler{
		sessionUseCase: sessionUseCase,
		sanitizer:      bluemonday.UGCPolicy(),
		hub: hub,
	}

	e.GET("/ws", handler.ws)
}

func (h *handler) ws(ctx echo.Context) error {
	conn, err := chat.Upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	client := &chat.Client{Hub: h.hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
	return nil
}
