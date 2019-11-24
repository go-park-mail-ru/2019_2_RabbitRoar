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
	Host    *Player
	Players []Player
	State   State
	Model   models.Game
	EvChan	chan EventWrapper
	logger  logging.Logger
}

func (game *Game) Run() {
	game.logger.Info("Starting game loop.")
	for {
		ew := <- game.EvChan

		if ew.Event.Type == WsRun {
			var allPlayersInfo []PlayerInfo

			for _, p := range game.Players {
				allPlayersInfo = append(allPlayersInfo, p.Info)
			}

			for _, p := range game.Players {
				noticeEvent := NewEvent(UserConnected, game.Model.Name, game.Model.PackName, allPlayersInfo)
				p.Conn.GetSendChan() <- *noticeEvent
			}
		}

		game.State = game.State.Handle(ew)
	}
}