package game

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type UseCase interface {
	Create(g *models.Game, userID int) error
	Fetch(page int) (*[]models.Game, error)
	JoinPlayerToGame(u models.User, gameID uuid.UUID) (*models.Game, error)
	KickPlayerFromGame(playerID int) error
	GetGameIDByUserID(userID int) (uuid.UUID, error)

	NewConnectionWrapper(ws *websocket.Conn) ConnectionWrapper
	JoinConnectionToGame(gameID uuid.UUID, userID int, conn ConnectionWrapper) error
}
