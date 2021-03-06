package game

import "math/rand"

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

	payload, ok := e.Event.Payload.(map[string]interface{})
	if !ok {
		s.Game.logger.Info("PendAnswer: invalid payload, keep state.")
		return s
	}

	playerAnswer, ok := payload["answer"].(string)
	if !ok {
		s.Game.logger.Info("PendAnswer: invalid payload answer, keep state.")
		return s
	}

	var ev Event

	ev = Event{
		Type: AnswerGivenBack,
		Payload: AnswerGivenBackPayload{
			PlayerAnswer: playerAnswer,
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

	// MOCK BLOCK REMOVE ME
	nextStateMock := &PendQuestionChoose{
		BaseState: BaseState{Game: s.Game},
	}

	randIdx := rand.Int() % len(s.Game.Players)
	nextStateMock.respondentID = s.Game.Players[randIdx].Info.ID

	ev = Event{
		Type: RequestQuestion,
		Payload: RequestQuestionPayload{
			PlayerID: nextStateMock.respondentID,
		},
	}
	s.Game.BroadcastEvent(ev)
	return nextStateMock
	// MOCK BLOCK REMOVE ME

	s.Game.logger.Info("PendAnswer: moving to the next state %v.", nextState)
	return nextState
}
