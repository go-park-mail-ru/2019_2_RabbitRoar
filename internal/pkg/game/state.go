package game

type StateType int

const (
	Pending StateType = iota
	Running StateType = iota
)

type State interface {
	GetType() StateType
	Handle(e Event) State
}
