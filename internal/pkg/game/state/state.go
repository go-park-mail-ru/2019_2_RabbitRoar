package state

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"

type PendPlayers struct {
}

func (state *PendPlayers) Handle(e game.Event) game.State {
	return nil
}
