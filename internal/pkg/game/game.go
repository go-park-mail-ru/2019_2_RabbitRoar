package game

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"reflect"
)

type Player struct {
	Info PlayerInfo
	Conn ConnectionWrapper
}

type Game struct {
	Host    *Player
	Players []Player
	State   State
	Model   models.Game
	logger  logging.Logger
}

func (game *Game) Run() {
	game.logger.Info("Starting game loop.")
	for {
		game.logger.Info("Getting event")
		e, err := game.getEvent()
		if err != nil {
			game.logger.Error(err)
			continue
		}

		game.State = game.State.Handle(e)
	}
}

func (game *Game) getEvent() (EventWrapper, error) {
	game.logger.Info("Getting event")

	var sc []reflect.SelectCase

	sc = append(sc, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(game.Host.Conn.GetReceiveChan()),
	})

	for _, p := range game.Players {
		game.logger.Info("Adding to SelectCase: ", p.Info.Username)
		sc = append(sc, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(p.Conn.GetReceiveChan()),
		})
	}

	_, eventWrap, _ := reflect.Select(sc)
	if eventWrap.IsValid() {
		return eventWrap.Interface().(EventWrapper), nil
	}

	return EventWrapper{}, errors.New("invalid event received")
}
