package game

type State interface {
	Handle(e EventWrapper) State
}

type BaseState struct {
	Game *Game
}
