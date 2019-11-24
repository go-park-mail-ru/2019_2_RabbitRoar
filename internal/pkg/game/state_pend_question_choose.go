package game

type PendQuestionChoose struct {
	BaseState
}

func (s *PendQuestionChoose) GetType() StateType {
	return Running
}

func (s *PendQuestionChoose) Handle(e EventWrapper) State {
	log.Info("PendQustionChosen: got event: ", e)
	return s
}
