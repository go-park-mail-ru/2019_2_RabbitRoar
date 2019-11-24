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

func (state *PendPlayers) Handle(ew game.EventWrapper) game.State {
	log.Infof("Handler received event: ", ew)

	return state
}
