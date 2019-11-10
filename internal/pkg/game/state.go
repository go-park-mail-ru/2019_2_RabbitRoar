package game

type State interface {
	handle(Event) *State
}
