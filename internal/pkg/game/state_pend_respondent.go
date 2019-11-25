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

	if e.Event.Type != RespondentReady {
		s.Game.logger.Info(
			"PendRespondent: got unexpected event %s, expected %s. ",
			e.Event.Type,
			RespondentReady,
		)
		return s
	}

	ev := Event{
		Type:    RequestAnswer,
		Payload: RequestAnswerPayload{
			PlayerID:e.SenderID,
		},
	}
	s.Game.BroadcastEvent(ev)

	nextState := &PendAnswer{
		BaseState:  BaseState{Game: s.Game},
		PlayerID:   e.SenderID,
		ThemeID:    s.ThemeID,
		QuestionID: s.QuestionID,
	}

	return nextState
}
