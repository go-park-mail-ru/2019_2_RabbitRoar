package game

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type PendVerdictState struct {
	BaseState
	stopTimer *time.Timer
}

func NewPendVerdictState(g *Game, ctx *StateContext) State {
	answer := g.Questions.GetAnswer(ctx.ThemeIdx, ctx.QuestionIdx)
	e := Event{
		Type: RequestVerdict,
		Payload: RequestVerdictPayload{
			CorrectAnswer: answer,
		},
	}
	g.Notify(e, g.Host)

	return &PendVerdictState{
		BaseState: BaseState{
			Game: g,
			Ctx:  ctx,
		},
		stopTimer: time.NewTimer(
			viper.GetDuration("internal.pend_verdict_duration") * time.Second,
		),
	}
}

func (s *PendVerdictState) Handle(ew EventWrapper) State {
	s.Game.logger.Info("PendVerdict: got event: ", ew)

	var nextState State

	select {
	case t := <-s.stopTimer.C:
		s.Game.logger.Info("PendVerdict: pending time exceeded: ", t.String())

		s.onVerdictCorrect()

		s.Ctx.QuestionSelectorID = s.Ctx.RespondentID
		nextState = NewPendQuestionChosenState(s.Game, s.Ctx)

	default:
		if err := s.validateEvent(ew); err != nil {
			s.Game.logger.Info(err)
			return s
		}

		switch {
		case ew.Event.Type == VerdictCorrect:
			s.onVerdictCorrect()

			s.Ctx.QuestionSelectorID = s.Ctx.RespondentID
			nextState = NewPendQuestionChosenState(s.Game, s.Ctx)

		case ew.Event.Type == VerdictWrong:
			s.onVerdictWrong()

			nextState = NewPendRespondentState(s.Game, s.Ctx)

		default:
			nextState = NewGameEndedState(s.Game, s.Ctx)
		}
	}

	s.Game.logger.Info("PendVerdict: moving to the next state %v.", nextState)
	return nextState
}

func (s *PendVerdictState) validateEvent(ew EventWrapper) error {
	if ew.SenderID != s.Game.Host.Info.ID {
		return errors.New(
			fmt.Sprintf(
				"PendVerdict: got event from unexpected player %d, expected %d. ",
				ew.SenderID,
				s.Game.Host.Info.ID,
			),
		)
	}

	if ew.Event.Type != VerdictCorrect && ew.Event.Type != VerdictWrong {
		return errors.New(
			fmt.Sprintf(
				"PendVerdict: got unexpected event %s, expected %s or %s. ",
				ew.Event.Type,
				VerdictCorrect,
				VerdictWrong,
			),
		)
	}

	return nil
}

func (s *PendVerdictState) notifyAllPlayersOfVerdict(verdict bool, correctAnswer string) {
	e := Event{
		Type:    VerdictGivenBack,
		Payload: VerdictPayload{
			Verdict:       verdict,
			CorrectAnswer: correctAnswer,
		},
	}
	s.Game.BroadcastEvent(e)
}

func (s *PendVerdictState) onVerdictCorrect() {
	questionCost := (s.Ctx.QuestionIdx + 1) * 100
	s.Game.UpdatePlayerScore(s.Ctx.RespondentID, questionCost)

	correctAnswer := s.Game.Questions.GetAnswer(s.Ctx.ThemeIdx, s.Ctx.QuestionIdx)
	s.notifyAllPlayersOfVerdict(true, correctAnswer)
}

func (s *PendVerdictState) onVerdictWrong() {
	questionCost := (s.Ctx.QuestionIdx + 1) * 100
	s.Game.UpdatePlayerScore(s.Ctx.RespondentID, -questionCost)

	s.notifyAllPlayersOfVerdict(false, "")
}