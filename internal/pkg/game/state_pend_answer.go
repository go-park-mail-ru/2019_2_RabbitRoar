package game

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type PendAnswerState struct {
	BaseState
}

func NewPendAnswerState(g *Game, ctx *StateContext) State {
	e := Event{
		Type:    RequestAnswer,
		Payload: RequestAnswerPayload{
			RespondentID: ctx.RespondentID,
		},
	}
	g.BroadcastEvent(e)

	g.StopTimer = time.NewTimer(
		viper.GetDuration("internal.pend_answer_duration") * time.Second,
	)

	return &PendAnswerState{
		BaseState:  BaseState{
			Game: g,
			Ctx:  ctx,
		},
	}
}

func (s *PendAnswerState) Handle(ew EventWrapper) State {
	s.Game.logger.Info("PendAnswer: got event: ", ew)

	var nextState State

	switch ew.Event.Type {
	case PendingExceeded:
		s.Game.logger.Info("PendAnswer: pending time exceeded")

		questionCost := (s.Ctx.QuestionIdx + 1) * 100
		s.Game.UpdatePlayerScore(s.Ctx.RespondentID, -questionCost)

		s.notifyAllPlayersOfNoGivenAnswer()

		nextState = NewPendRespondentState(s.Game, s.Ctx)

	default:
		if err := s.validateEvent(ew); err != nil {
			s.Game.logger.Info(err)
			return s
		}

		givenAnswer, err := s.getPlayerGivenAnswer(ew)
		if err != nil {
			s.Game.logger.Info(err)
			return s
		}

		s.notifyAllPlayersOfGivenAnswer(givenAnswer, ew.SenderID)

		nextState = NewPendVerdictState(s.Game, s.Ctx)
	}

	s.Game.logger.Info("PendAnswer: moving to the next state %v.", nextState)
	return nextState
}

func (s *PendAnswerState) validateEvent(ew EventWrapper) error {
	if ew.SenderID != s.Ctx.RespondentID {
		return errors.New(
			fmt.Sprintf(
				"PendAnswer: got event from unexpected player %d, expected %d. ",
				ew.SenderID,
				s.Ctx.RespondentID,
			),
		)
	}

	if ew.Event.Type != AnswerGiven {
		return errors.New(
			fmt.Sprintf(
				"PendAnswer: got unexpected event %s, expected %s. ",
				ew.Event.Type,
				AnswerGiven,
			),
		)
	}

	return nil
}

func (s *PendAnswerState) getPlayerGivenAnswer(ew EventWrapper) (string, error) {
	payload, ok := ew.Event.Payload.(map[string]interface{})
	if !ok {
		return "", errors.New("PendAnswer: invalid payload, keep state.")
	}

	playerAnswer, ok := payload["answer"].(string)
	if !ok {
		return "", errors.New("PendAnswer: invalid payload answer, keep state.")
	}

	return playerAnswer, nil
}

func (s *PendAnswerState) notifyAllPlayersOfGivenAnswer(answer string, respondentID int) {
	e := Event{
		Type: AnswerGivenBack,
		Payload: AnswerGivenBackPayload{
			PlayerAnswer: answer,
			RespondentID: respondentID,
		},
	}
	s.Game.BroadcastEvent(e)
}

func (s *PendAnswerState) notifyAllPlayersOfNoGivenAnswer() {
	e := Event{
		Type: VerdictGivenBack,
		Payload: VerdictPayload{
			Verdict:       false,
			CorrectAnswer: "",
			Players:       s.Game.GatherPlayersInfo(),
		},
	}
	s.Game.BroadcastEvent(e)

	time.Sleep(5 * time.Second)
}