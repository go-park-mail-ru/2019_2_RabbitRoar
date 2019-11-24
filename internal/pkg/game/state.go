package game

type StateType int

const (
	Pending StateType = iota
	Running StateType = iota
)

type State interface {
	GetType() StateType
	//SetGame(g *Game)
	Handle(e EventWrapper) State
}

type BaseState struct {
	Game *Game
}
