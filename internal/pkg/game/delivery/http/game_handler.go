package http

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type handler struct {
	usecase game.UseCase
}

var upgrader = websocket.Upgrader{}

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
	// group.POST("/:uuid/join")
	// group.DELETE("/leave")
	group.GET("/ws", handler.ws)
}

func (gh *handler) self(ctx echo.Context) error {
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "cannot parse page",
			Internal: err,
		}
	}

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

	return ctx.JSON(http.StatusOK, content)
}

func (gh *handler) create(ctx echo.Context) error {

}

func (gh *handler) ws(ctx echo.Context) error {
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	var wg sync.WaitGroup

	wg.Add(2)

	go func(readChan chan []byte, stop chan bool) {
		defer wg.Done()

		for {
			select {
			case <-stop:
				close(readChan)
				log.Info("Stopped reading from websocket manually")
				return

			default:
				_, msg, err := ws.ReadMessage()
				if err != nil {
					close(readChan)
					log.Info(err)
					return
				}

				readChan <- msg
			}
		}
	}()

	go func(writeChan chan []byte, stop chan bool) {
		defer wg.Done()

		ticker := time.NewTicker(10 * time.Second)

		for {
			select {
			case <-stop:
				log.Info("Stopped writing into websocket manually")
				return

			case msg := <-writeChan:
				err = ws.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Info(err)
					return
				}

			case <-ticker.C:
				err := ws.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					log.Info(err)
					return
				}
			}
		}
	}()

	wg.Wait()

	return nil
}
