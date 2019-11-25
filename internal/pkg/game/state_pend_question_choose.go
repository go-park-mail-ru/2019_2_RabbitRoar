package game

type PendQuestionChoose struct {
	BaseState
}

func (s *PendQuestionChoose) Handle(e EventWrapper) State {
	log.Info("PendQustionChosen: got event: ", e)
	return s
}
