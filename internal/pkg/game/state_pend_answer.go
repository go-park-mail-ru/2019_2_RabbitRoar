package game

type PendAnswer struct {
	BaseState
	PlayerID   int
	ThemeID    int
	QuestionID int
}

func (s *PendAnswer) getAnswer(themeID, questionID int) string {
	themeSlice := s.Game.Questions.([]interface{})

	for themeidx, theme := range themeSlice {
		theme := theme.(map[string]interface{})
		if themeID == themeidx {
			questionSlice := theme["questions"].([]interface{})
			for questionidx, question := range questionSlice {
				question := question.(map[string]interface{})
				if questionidx == questionID {
					return question["answer"].(string)
				}
			}
		}
	}

	return ""
}

func (s *PendAnswer) Handle(e EventWrapper) State {
	s.Game.logger.Info("PendAnswer: got event: ", e)

	if e.Event.Type != AnswerGiven {
		s.Game.logger.Info(
			"PendAnswer: got unexpected event %s, expected %s. ",
			e.Event.Type,
			AnswerGiven,
		)
		return s
	}

	payload, ok := e.Event.Payload.(AnswerGivenPayload)
	if !ok {
		s.Game.logger.Info("PendAnswer: invalid payload, keep state.")
		return s
	}

	var ev Event

	ev = Event{
		Type: AnswerGivenBack,
		Payload: AnswerGivenBackPayload{
			PlayerAnswer: payload.Answer,
			PlayerID:     e.SenderID,
		},
	}
	s.Game.BroadcastEvent(ev)

	answer := s.getAnswer(s.ThemeID, s.QuestionID)
	ev = Event{
		Type: RequestVerdict,
		Payload: RequestVerdictPayload{
			HostID: s.Game.Host.Info.ID,
			Answer: answer,
		},
	}
	s.Game.Host.Conn.GetSendChan() <- ev

	nextState := &PendVerdict{
		BaseState: BaseState{Game: s.Game},
		PlayerID:  e.SenderID,
	}

	s.Game.logger.Info("PendAnswer: moving to the next state %v.", nextState)
	return nextState
}
