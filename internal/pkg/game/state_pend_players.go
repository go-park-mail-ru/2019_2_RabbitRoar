package game

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("State")

type PendPlayers struct {
	BaseState
}

func (s *PendPlayers) GetType() StateType {
	return Pending
}

func (s *PendPlayers) Handle(e EventWrapper) State {
	log.Info("PendPlayers: got event: ", e)
	if e.Event.Type == PlayerReadyFront {
		var playersReady int

		for idx, pl := range s.Game.Players {
			if pl.Info.ID == e.SenderID {
				s.Game.Players[idx].Info.Ready = true
			}

			if s.Game.Players[idx].Info.Ready {
				playersReady++
			}
		}


		var players = make([]PlayerInfo, 0, len(s.Game.Players))
		for _, pl := range s.Game.Players {
			players = append(players, pl.Info)
		}

		ev := Event{
			Type:    PlayerReadyBack,
			Payload: players,
		}

		for _, pl := range s.Game.Players {
			pl.Conn.GetSendChan() <- ev
		}

		if playersReady == s.Game.Model.PlayersCapacity {
			return &PendQuestionChoose{BaseState{Game:s.Game}}
		}
	}
	return s
}

type PendQuestionChoose struct {
	BaseState
}

func (s *PendQuestionChoose) GetType() StateType {
	return Running
}

func (s *PendQuestionChoose) Handle(e EventWrapper) State {
	log.Info("PendQustionChosen: got event: ", e)
	return s
}
