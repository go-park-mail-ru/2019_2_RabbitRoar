package game

type State interface {
	Handle(e Event) State
}
