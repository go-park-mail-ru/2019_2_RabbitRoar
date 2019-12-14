package game

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/op/go-logging"
	"math/rand"
)

type Player struct {
	Info PlayerInfo
	Conn ConnectionWrapper
}

type Game struct {
	Host      *Player
	Players   []Player
	State     State
	Model     models.Game
	Questions *QuestionTable
	EvChan    chan EventWrapper
	Started   bool
	logger    logging.Logger
}

func (g *Game) Run(killChan chan uuid.UUID) {
	defer func() {
		g.logger.Info("Started closing connections")
		for _, p := range g.Players {
			g.logger.Infof("Trying to close player connection: %d", p.Info.ID)
			if p.Conn.IsRunning() {
				g.logger.Info("Connection is running. Stopping connection")
				p.Conn.Stop()
				g.logger.Info("Connection stopped")
			}
		}
		g.logger.Infof("All connections stopped. Game is ready to be deleted. UUID: %s", g.Model.UUID.String())
		killChan <- g.Model.UUID
	}()

	g.logger.Info("Starting game loop.")
	g.State = NewPendPlayersState(g)

	for {
		if len(g.Players) == 0 {
			return
		}
		g.logger.Info("Pending event...")
		ew := <-g.EvChan
		g.logger.Info("Got event: ", ew)

		if ew.Event.Type == WsUpdated {
			var allPlayersInfo []PlayerInfo

			for _, p := range g.Players {
				allPlayersInfo = append(allPlayersInfo, p.Info)
			}

			noticeEvent := Event{
				Type: UserConnected,
				Payload: UserConnectedPayload{
					RoomName: g.Model.Name,
					PackName: g.Model.PackName,
					Host:     g.Host.Info,
					Players:  allPlayersInfo,
				},
			}

			g.BroadcastEvent(noticeEvent)
			continue
		}

		g.State = g.State.Handle(ew)
		if g.State == nil {
			return
		}
	}
}

func (g *Game) Notify(e Event, player *Player) {
	if player.Conn.IsRunning(){
		player.Conn.GetSendChan() <- e
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

func (g *Game) GatherPlayersInfo() []PlayerInfo {
	playersInfo := make([]PlayerInfo, 0, len(g.Players))

	for _, pl := range g.Players {
		playersInfo = append(playersInfo, pl.Info)
	}

	return playersInfo
}

func (g *Game) GetRandPlayerID() int {
	playerID := g.Host.Info.ID

	for playerID == g.Host.Info.ID {
		randIdx := rand.Int() % len(g.Players)
		playerID = g.Players[randIdx].Info.ID
	}

	return playerID
}

func (g *Game) GetNextPlayerID(playerID int) int {
	prevPlayerIdx, err := g.getPlayerIdxByPlayerID(playerID)
	if err != nil {
		return g.GetRandPlayerID()
	}

	nextPlayerIdx := (prevPlayerIdx + 1) % len(g.Players)

	for nextPlayerIdx == g.Host.Info.ID {
		nextPlayerIdx = (nextPlayerIdx + 1) % len(g.Players)
	}

	return nextPlayerIdx
}

func (g *Game) UpdatePlayerScore(playerID, score int) {
	playerIdx, err := g.getPlayerIdxByPlayerID(playerID)
	if err != nil {
		return
	}

	g.Players[playerIdx].Info.Score += score
}

func (g *Game) getPlayerIdxByPlayerID(playerID int) (int, error) {
	for idx, p := range g.Players {
		if p.Info.ID == playerID {
			return idx, nil
		}
	}

	return 0, nil
}