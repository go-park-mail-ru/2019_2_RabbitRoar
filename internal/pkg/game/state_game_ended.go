package game

import (
	"github.com/spf13/viper"
	"time"
)

type GameEndedState struct {
	BaseState
}

func NewGameEndedState(g *Game, ctx *StateContext) State {
	e := Event{
		Type:    GameEnded,
		Payload: GameEndedPayload{
			Players: g.GatherPlayersInfo(),
		},
	}

	g.BroadcastEvent(e)

	addScoreForPlayers(g)

	g.StopTimer = time.NewTimer(
		viper.GetDuration("internal.pend_game_ended_duration") * time.Second,
	)
	
	return &GameEndedState{
		BaseState: BaseState{
			Game: g,
			Ctx:  ctx,
		},
	}
}

func (s *GameEndedState) Handle(ew EventWrapper) State {
	s.Game.logger.Info("GameEnded: got event: ", ew)

	for {
		switch ew.Event.Type {
		case PendingExceeded:
			s.Game.logger.Info("GameEnded: stopping game.")
			return nil

		default:
			// Handle voting here
		}
	}
}

func addScoreForPlayers(g *Game) {
	for _, p := range g.Players {
		g.UpdateUserRating(p.Info)
	}
}