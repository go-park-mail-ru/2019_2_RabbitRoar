package game

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type PendPlayersState struct {
	BaseState
	stopTimer *time.Timer
}

func NewPendPlayersState(g *Game) State {
	return &PendPlayersState{
		BaseState: BaseState{
			Game: g,
			Ctx: &StateContext{
				QuestionSelectorID: 0,
				ThemeIdx:           0,
				QuestionIdx:        0,
				RespondentID:       0,
			},
		},
		stopTimer: time.NewTimer(
			viper.GetDuration("internal.pend_players_duration") * time.Second,
		),
	}
}

func (s *PendPlayersState) Handle(ew EventWrapper) State {
	s.Game.logger.Info("PendPlayers: got event: ", ew)

	select {
	case t := <-s.stopTimer.C:
		s.Game.logger.Info("PendPlayers: pending time exceeded: ", t.String())
		return nil

	default:
		if err := s.validateEvent(ew); err != nil {
			s.Game.logger.Info(err)
			return s
		}

		playersReady := s.updateReadyPlayers(ew)
		s.sendPlayersInfo()

		if playersReady != s.Game.Model.PlayersCapacity {
			s.Game.logger.Info(
				"PendPlayers: players ready %d/%d, keep state.",
				playersReady,
				s.Game.Model.PlayersCapacity,
			)
			return s
		}

		s.Game.Started = true

		s.Ctx.QuestionSelectorID = s.Game.GetRandPlayerID()
		nextState := NewPendQuestionChosenState(s.Game, s.Ctx)

		s.Game.logger.Info("PendPlayers: moving to the next state %v.", nextState)

		return nextState
	}
}

func (s *PendPlayersState) validateEvent(ew EventWrapper) error {
	if ew.Event.Type != PlayerReadyFront {
		return errors.New(
			fmt.Sprintf(
				"PendPlayers: got unexpected event %s, expected %s. ",
				ew.Event.Type,
				PlayerReadyFront,
			),
		)
	}

	return nil
}

func (s *PendPlayersState) sendPlayersInfo() {
	e := Event{
		Type:    PlayerReadyBack,
		Payload: PlayerReadyBackPayload{
			Players: s.Game.GatherPlayersInfo(),
		},
	}

	s.Game.BroadcastEvent(e)
}

func (s *PendPlayersState) updateReadyPlayers(ew EventWrapper) int {
	var playersReady int
	for playerID, p := range s.Game.Players {
		if p.Info.ID == ew.SenderID {
			s.Game.Players[playerID].Info.Ready = !p.Info.Ready
		}

		if s.Game.Players[playerID].Info.Ready {
			playersReady++
		}
	}

	return playersReady
}
