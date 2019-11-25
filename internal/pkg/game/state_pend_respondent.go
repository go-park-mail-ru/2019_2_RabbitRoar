package game

type PendRespondent struct {
	BaseState
	ThemeID    int
	QuestionID int
}

func (s *PendRespondent) GetType() StateType {
	return Running
}

func (s *PendRespondent) Handle(e EventWrapper) State {
	s.Game.logger.Info("PendRespondent: got event: ", e)
	return s
}
