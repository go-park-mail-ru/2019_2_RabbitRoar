package state

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("State")

type BasicState struct {
	game *game.Game
}

type PendPlayers struct {
	BasicState
}

func (s *PendPlayers) GetType() game.StateType {
	return game.Pending
}

func (s *PendPlayers) Handle(e game.EventWrapper) game.State {
	log.Info("PendPlayers: got event: ", e)
	if e.Event.Type == game.PlayerReadyFront {
		var playersReady int

		for _, pl := range s.game.Players {
			if pl.Info.ID == e.SenderID {
				pl.Info.Ready = true
				return s
			}

			if pl.Info.Ready {
				playersReady++
			}
		}


		var players = make([]game.PlayerInfo, len(s.game.Players))
		for _, pl := range s.game.Players {
			players = append(players, pl.Info)
		}

		ev := game.Event{
			Type:    game.PlayerReadyBack,
			Payload: players,
		}

		for _, pl := range s.game.Players {
			pl.Conn.GetSendChan() <- ev
		}

		if playersReady == s.game.Model.PlayersCapacity {
			return &PendQuestionChoose{BasicState{game:s.game}}
		}
	}
	return s
}

type PendQuestionChoose struct {
	BasicState
}

func (s *PendQuestionChoose) GetType() game.StateType {
	return game.Running
}

func (s *PendQuestionChoose) Handle(e game.EventWrapper) game.State {
	log.Info("PendQustionChosen: got event: ", e)
	return s
}
