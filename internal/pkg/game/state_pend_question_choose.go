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

	if e.Event.Type != QuestionChosen {
		s.Game.logger.Infof(
			"PendQuestionChosen: got unexpected event %s, expected %s.",
			e.Event.Type,
			QuestionChosen,
		)
		return s
	}

	payload, ok := e.Event.Payload.(map[string]interface{})
	if !ok {
		s.Game.logger.Info("PendQuestionChosen: got invalid payload, keep old state.")
	}

	questionIdxFloat, ok := payload["question_idx"].(float64)
	if !ok {
		s.Game.logger.Info("PendQuestionChosen: got invalid payload, keep old state.")
	}

	themeIdxFloat, ok := payload["theme_idx"].(float64)
	if !ok {
		s.Game.logger.Info("PendQuestionChosen: got invalid payload, keep old state.")
	}

	themeIdx := int(themeIdxFloat)
	questionIdx := int(questionIdxFloat)

	if themeIdx < 0 || themeIdx > 4 {
		s.Game.logger.Info("PendQuestionChosen: got invalid theme coords, keep old state.")
		return s
	}

	if questionIdx < 0 || questionIdx > 4 {
		s.Game.logger.Info("PendQuestionChosen: got invalid question coords, keep old state.")
		return s
	}

	ev := Event{
		Type: RequestRespondent,
		Payload: RequestRespondentPayload{
			Question:   s.getQuestion(themeIdx, questionIdx),
			ThemeID:    themeIdx,
			QuestionID: questionIdx,
		},
	}

	s.Game.BroadcastEvent(ev)

	nextState := &PendRespondent{
		BaseState:  BaseState{Game: s.Game},
		ThemeID:    themeIdx,
		QuestionID: questionIdx,
	}

	s.Game.logger.Info("PendQuestionChoose: moving to the next state %v.", nextState)
	return nextState
}
