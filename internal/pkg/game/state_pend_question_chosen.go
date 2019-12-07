package game

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type PendQuestionChoose struct {
	BaseState
	stopTimer *time.Timer
}

func NewPendQuestionChosenState(g *Game, ctx *StateContext) State {
	e := Event{
		Type:    GameStart,
		Payload: GameStartPayload{
			Themes: g.Questions.GetThemes(),
		},
	}
	g.BroadcastEvent(e)

	e = Event{
		Type: RequestQuestion,
		Payload: RequestQuestionPayload{
			QuestionSelectorID: ctx.QuestionSelectorID,
		},
	}
	g.BroadcastEvent(e)

	return &PendQuestionChoose{
		BaseState: BaseState{
			Game: g,
			Ctx:  ctx,
		},
		stopTimer: time.NewTimer(
			viper.GetDuration("internal.pend_question_duration") * time.Second,
		),
	}
}

func (s *PendQuestionChoose) Handle(ew EventWrapper) State {
	s.Game.logger.Info("PendQuestionChosen: got event: ", ew)

	var nextState State

	select {
	case t := <-s.stopTimer.C:
		s.Game.logger.Info("PendQuestionChosen: pending time exceeded: ", t.String())

		themeIdx, questionIdx, err := s.Game.Questions.GetRandAvailableQuestionIndexes()
		if err != nil {
			s.Game.logger.Info(err)
			return nil
		}

		s.Ctx.ThemeIdx = themeIdx
		s.Ctx.QuestionIdx = questionIdx
		nextState = NewPendRespondentState(s.Game, s.Ctx)

	default:
		if err := s.validateEvent(ew); err != nil {
			s.Game.logger.Info(err)
			return s
		}

		themeIdx, questionIdx, err := s.getQuestionIndexes(ew)
		if err != nil {
			s.Game.logger.Info(err)
		}

		s.Ctx.ThemeIdx = themeIdx
		s.Ctx.QuestionIdx = questionIdx
		nextState = NewPendRespondentState(s.Game, s.Ctx)
	}

	s.Game.logger.Info("PendQuestionChoose: moving to the next state %v.", nextState)
	return nextState
}


func (s *PendQuestionChoose) validateEvent(ew EventWrapper) error {
	if ew.SenderID != s.Ctx.QuestionSelectorID {
		return errors.New(
			fmt.Sprintf(
				"PendQuestionChosen: got event from unexpected player %d, expected %d.",
				ew.SenderID,
				s.Ctx.QuestionSelectorID,
			),
		)
	}

	if ew.Event.Type != QuestionChosen {
		return errors.New(
			fmt.Sprintf(
				"PendQuestionChosen: got unexpected event %s, expected %s.",
				ew.Event.Type,
				QuestionChosen,
			),
		)
	}

	return nil
}

func (s *PendQuestionChoose) getQuestionIndexes(ew EventWrapper) (int, int, error) {
	payload, ok := ew.Event.Payload.(map[string]interface{})
	if !ok {
		return 0, 0, errors.New("PendQuestion: got invalid payload, keep old state.")
	}

	questionIdxFloat, ok := payload["question_idx"].(float64)
	if !ok {
		return 0, 0, errors.New("PendQuestion: got invalid payload, keep old state.")
	}

	themeIdxFloat, ok := payload["theme_idx"].(float64)
	if !ok {
		return 0, 0, errors.New("PendQuestion: got invalid payload, keep old state.")
	}

	themeIdx := int(themeIdxFloat)
	questionIdx := int(questionIdxFloat)

	if themeIdx < 0 || themeIdx > 4 {
		return 0, 0, errors.New("PendQuestion: got invalid theme coords, keep old state.")
	}

	if questionIdx < 0 || questionIdx > 4 {
		return 0, 0, errors.New("PendQuestion: got invalid question coords, keep old state.")
	}

	return themeIdx, questionIdx, nil
}