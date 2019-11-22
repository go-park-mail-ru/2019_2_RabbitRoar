package game

import (
	"github.com/op/go-logging"
	"reflect"
)

var log = logging.MustGetLogger("game")

type Player struct {
	Info PlayerInfo
	Conn ConnectionWrapper
}

type Game struct {
	Host    Player
	Players []Player
	State   State
	logger logging.Logger
}

func (game *Game) Run() {
	game.logger.Info("Starting game loop.")
	for {
		game.logger.Info("Getting event")
		game.getEvent()
	}
}

func (game *Game) getEvent() EventWrapper {
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
	eventWrap.IsValid() //TODO: handle it
	return eventWrap.Interface().(EventWrapper)
}
