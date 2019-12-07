package game

type State interface {
	Handle(e EventWrapper) State
}

type StateContext struct {
	QuestionSelectorID int
	ThemeIdx           int
	QuestionIdx        int
	RespondentID       int
}

type BaseState struct {
	Game *Game
	Ctx  *StateContext
}
