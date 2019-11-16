package game

import "reflect"

type Player struct {
	Info PlayerInfo
	Conn Connection
}

type Game struct {
	Host    Player
	Players []Player
	State   State
}

func (game *Game) Run() {
	// for {

	// }
}

func (game *Game) getEvent() {
	var sc []reflect.SelectCase

	sc = append(sc, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(game.Host.Conn.GetReceiveChan()),
	})

	for _, p := range game.Players {
		sc = append(sc, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(p.Conn.GetReceiveChan()),
		})
	}

	_, eventWrap, _ := reflect.Select(sc)
}
