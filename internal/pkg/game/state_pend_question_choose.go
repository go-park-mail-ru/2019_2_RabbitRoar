package game

type PendQuestionChoose struct {
	BaseState
	respondentID int
}

func (s *PendQuestionChoose) getQuestion(themeID, questionID int) string {
	themeSlice := s.Game.Questions.([]interface{})

	for themeidx, theme := range themeSlice {
		theme := theme.(map[string]interface{})
		if themeID == themeidx {
			questionSlice := theme["questions"].([]interface{})
			for questionidx, question := range questionSlice {
				question := question.(map[string]interface{})
				if questionidx == questionID {
					return question["text"].(string)
				}
			}
		}
	}

	return ""
}

func (s *PendQuestionChoose) GetType() StateType {
	return Running
}

func (s *PendQuestionChoose) Handle(e EventWrapper) State {
	s.Game.logger.Info("PendQuestionChosen: got event: ", e)

	if e.SenderID != s.respondentID {
		s.Game.logger.Infof(
			"PendQuestionChosen: got event from unexpected player %d, expected %d.",
			e.SenderID,
			s.respondentID,
		)
		return s
	}

	if e.Event.Type == QuestionChosen {
		s.Game.logger.Infof(
			"PendQuestionChosen: got unexpected event %s, expected %s.",
			e.Event.Type,
			QuestionChosen,
		)
		return s
	}

	payload, ok := e.Event.Payload.(QuestionChosenPayload)
	if !ok {
		s.Game.logger.Info("PendQuestionChosen: got invalid payload, keep old state.")
		return s
	}

	if payload.ThemeIdx < 0 || payload.ThemeIdx > 4 {
		s.Game.logger.Info("PendQuestionChosen: got invalid theme coords, keep old state.")
		return s
	}

	if payload.QuestionIdx < 0 || payload.QuestionIdx > 4 {
		s.Game.logger.Info("PendQuestionChosen: got invalid question coords, keep old state.")
		return s
	}

	ev := Event{
		Type: RequestRespondent,
		Payload: RequestRespondentPayload{
			Question:   s.getQuestion(payload.ThemeIdx, payload.QuestionIdx),
			ThemeID:    payload.ThemeIdx,
			QuestionID: payload.QuestionIdx,
		},
	}

	s.Game.BroadcastEvent(ev)

	nextState := &PendRespondent{
		BaseState:  BaseState{Game: s.Game},
		ThemeID:    payload.ThemeIdx,
		QuestionID: payload.QuestionIdx,
	}

	return nextState
}
