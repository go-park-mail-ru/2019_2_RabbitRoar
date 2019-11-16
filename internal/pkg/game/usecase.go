package game

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type UseCase interface {
	GetByID(uuid uuid.UUID) (*models.Game, error)
	GetGameIDByUserID(userID int) (uuid.UUID, error)
	SQLCreate(g models.Game, u models.User) error
	Fetch(page int) (*[]models.Game, error)
	JoinPlayerToGame(playerID int, gameID uuid.UUID) (*models.Game, error)
	KickPlayerFromGame(playerID int) error

	NewConnection(ws *websocket.Conn) Connection
	MemCreate(g models.Game, u models.User) error
	JoinConnectionToGame(gameID uuid.UUID, u models.User, conn Connection) error
}
