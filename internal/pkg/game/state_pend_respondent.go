package game

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type PendRespondentState struct {
	BaseState
	stopTimer  *time.Timer
}

func NewPendRespondentState(g *Game, ctx *StateContext) State {
	e := Event{
		Type: RequestRespondent,
		Payload: RequestRespondentPayload{
			Question:    g.Questions.GetQuestion(ctx.ThemeIdx, ctx.QuestionIdx),
			ThemeIdx:    ctx.ThemeIdx,
			QuestionIdx: ctx.QuestionIdx,
		},
	}

	g.BroadcastEvent(e)

	return &PendRespondentState{
		BaseState:  BaseState{
			Game: g,
			Ctx:  ctx,
		},
		stopTimer:  time.NewTimer(
			viper.GetDuration("internal.pend_respondent_duration") * time.Second,
		),
	}
}

func (s *PendRespondentState) Handle(ew EventWrapper) State {
	s.Game.logger.Info("PendRespondent: got event: ", ew)

	var nextState State

	select {
	case t := <-s.stopTimer.C:
		s.Game.logger.Info("PendRespondent: pending time exceeded: ", t.String())

		s.Ctx.QuestionSelectorID = s.Game.GetNextPlayerID(s.Ctx.QuestionSelectorID)
		nextState = NewPendQuestionChosenState(s.Game, s.Ctx)

	default:
		if err := s.validateEvent(ew); err != nil {
			s.Game.logger.Info(err)
			return s
		}

		s.Ctx.RespondentID = ew.SenderID
		nextState = NewPendAnswerState(s.Game, s.Ctx)
	}

	s.Game.logger.Info("PendRespondent: moving to the next state %v.", nextState)
	return nextState
}


func (s *PendRespondentState) validateEvent(ew EventWrapper) error {
	if ew.SenderID == s.Game.Host.Info.ID {
		return errors.New(
			fmt.Sprintf(
				"PendRespondent: got event from unexpected player %s, expected any except %s. ",
				ew.Event.Type,
				s.Game.Host.Info.ID,
			),
		)
	}

	if ew.Event.Type != RespondentReady {
		return errors.New(
			fmt.Sprintf(
				"PendRespondent: got unexpected event %s, expected %s. ",
				ew.Event.Type,
				RespondentReady,
			),
		)
	}

	return nil
}