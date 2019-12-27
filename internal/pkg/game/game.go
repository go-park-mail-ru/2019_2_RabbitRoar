package game

import (
	"errors"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
	"github.com/google/uuid"
	"github.com/op/go-logging"
	"math/rand"
	"time"
)

type Player struct {
	Info PlayerInfo
	Conn ConnectionWrapper
}

type Game struct {
	Host      *Player
	Players   []*Player
	State     State
	Model     models.Game
	Questions *QuestionTable
	EvChan    chan EventWrapper
	Started   bool
	StopTimer *time.Timer
	UserRepo  user.Repository
	logger    logging.Logger
}

func (g *Game) Run(killChan chan uuid.UUID) {
	defer g.safeStop(killChan)

	g.logger.Info("Starting game loop.")
	g.State = NewPendPlayersState(g)

	for {
		if len(g.Players) == 0 {
			return
		}
		g.logger.Info("Pending event...")

		var ew EventWrapper

		select {
		case t := <- g.StopTimer.C:
			g.logger.Info("Current event pending time exceeded: ", t.String())
			ew = EventWrapper{
				SenderID: -1,
				Event:    &Event{
					Type:    PendingExceeded,
					Payload: &PendingExceededPayload{
						Time: t,
					},
				},
			}

		case ew = <-g.EvChan:
			g.logger.Info("Got event: ", ew)
		}

		if ew.Event.Type == WsUpdated {
			g.handleWSUpdated()
			continue
		}

		if ew.Event.Type == PlayerLeft {
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
	if len(g.Players) == 0 {
		return 0
	}

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

	for g.Players[nextPlayerIdx].Info.ID == g.Host.Info.ID {
		nextPlayerIdx = (nextPlayerIdx + 1) % len(g.Players)
	}

	return g.Players[nextPlayerIdx].Info.ID
}

func (g *Game) UpdatePlayerScore(playerID, score int) {
	playerIdx, err := g.getPlayerIdxByPlayerID(playerID)
	if err != nil {
		return
	}

	g.Players[playerIdx].Info.Score += score
}

func (g *Game) UpdateUserRating(playerInfo PlayerInfo) {
	u, err := g.UserRepo.GetByID(playerInfo.ID)
	if err != nil {
		g.logger.Info("Error finding user: ", err)
		return
	}

	g.logger.Infof("Changed user rating: %d + %d", u.Rating, playerInfo.Score)
	u.Rating += playerInfo.Score

	err = g.UserRepo.Update(*u)
	if err != nil {
		g.logger.Info("Error updating user: ", err)
		return
	}
}

func (g *Game) safeStop(killChan chan uuid.UUID) {
	g.logger.Info("Started closing connections")
	for _, p := range g.Players {
		g.logger.Info("Trying to close player connection. ID: ", p.Info.ID)
		if p.Conn.IsRunning() {
			g.logger.Info(" -- Connection is running. Stopping connection")
			p.Conn.Stop()
			g.logger.Info(" -- Connection stopped")
		}
	}
	g.logger.Info("All connections stopped. Game is ready to be deleted. UUID: ", g.Model.UUID.String())
	killChan <- g.Model.UUID
}

func (g *Game) handleWSUpdated() {
	var allPlayersInfo []PlayerInfo

	for _, p := range g.Players {
		allPlayersInfo = append(allPlayersInfo, p.Info)
	}

	noticeEvent := Event{
		Type: UserConnected,
		Payload: UserConnectedPayload{
			RoomName:        g.Model.Name,
			PackName:        g.Model.PackName,
			Host:            g.Host.Info,
			Players:         allPlayersInfo,
			QuestionsStatus: g.Questions.questionsAvailable,
		},
	}

	g.BroadcastEvent(noticeEvent)
}

func (g *Game) getPlayerIdxByPlayerID(playerID int) (int, error) {
	for idx, p := range g.Players {
		if p.Info.ID == playerID {
			return idx, nil
		}
	}

	return 0, errors.New("player with such ID is not found")
}