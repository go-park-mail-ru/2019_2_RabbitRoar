package game

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/op/go-logging"
)

type Player struct {
	Info PlayerInfo
	Conn ConnectionWrapper
}

type Game struct {
	Host      *Player
	Players   []Player
	Questions interface{}
	State     State
	Model     models.Game
	EvChan    chan EventWrapper
	logger    logging.Logger
}

func (g *Game) Run() {
	g.logger.Info("Starting game loop.")
	g.State = &PendPlayers{
		BaseState: BaseState{
			Game: g,
		},
	}

	for {
		ew := <- g.EvChan

		if ew.Event.Type == WsRun {
			var allPlayersInfo []PlayerInfo

			for _, p := range g.Players {
				allPlayersInfo = append(allPlayersInfo, p.Info)
			}

			noticeEvent := Event{
				Type:    UserConnected,
				Payload: UserConnectedPayload{
					RoomName: g.Model.Name,
					PackName: g.Model.PackName,
					Players:  allPlayersInfo,
				},
			}
			g.BroadcastEvent(noticeEvent)
		}

		g.State = g.State.Handle(ew)
		if g.State == nil {
			return
		}
	}
}

func (g *Game) BroadcastEvent(e Event) {
	for _, p := range g.Players {
		if !p.Conn.IsRunning() {
			continue
		}
		p.Conn.GetSendChan() <- e
	}
}
