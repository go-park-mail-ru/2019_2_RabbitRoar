package game

type PendVerdict struct {
	BaseState
	PlayerID   int
}

func (s *PendVerdict) GetType() StateType {
	return Running
}

func (s *PendVerdict) Handle(e EventWrapper) State {
	s.Game.logger.Info("PendVerdict: got event: ", e)
	return s
}

