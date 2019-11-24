package state

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/prometheus/common/log"
)

type BasicState struct {
	game *game.Game
}

type PendPlayers struct {
}

func (state *PendPlayers) GetType() game.StateType {
	return game.Pending
}

func (state *PendPlayers) Handle(e game.Event) game.State {
	log.Info()
	return state
}
