package game

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("State")

type PendPlayers struct {
	BaseState
}

func (s *PendPlayers) getThemes() [5]string {
	var themes [5]string
	themeSlice := s.Game.Questions.([]interface{})

	for i := 0; i < 5; i++ {
		theme := themeSlice[i].(map[string]interface{})
		themes[i] = theme["name"].(string)
	}

	return themes
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
				s.Game.Players[idx].Info.Ready = !s.Game.Players[idx].Info.Ready
			}

			if s.Game.Players[idx].Info.Ready {
				playersReady++
			}
		}

		// collect joined players
		var players = make([]PlayerInfo, 0, len(s.Game.Players))
		for _, pl := range s.Game.Players {
			players = append(players, pl.Info)
		}

		ev := Event{
			Type:    PlayerReadyBack,
			Payload: players,
		}

		s.Game.BroadcastEvent(ev)

		if playersReady == s.Game.Model.PlayersCapacity {
			ev = Event{
				Type:    GameStart,
				Payload: GameStartPayload{Themes:s.getThemes()},
			}
			s.Game.BroadcastEvent(ev)
			return &PendQuestionChoose{BaseState{Game:s.Game}}
		}
	}
	return s
}
