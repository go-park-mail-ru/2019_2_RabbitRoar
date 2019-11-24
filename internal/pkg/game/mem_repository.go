package game

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type MemRepository interface {
	Create(g *models.Game, host models.User) error
	Fetch(pageSize int, page int) (*[]models.Game, error)

	JoinPlayer(u models.User, gameID uuid.UUID) (*models.Game, error)
	JoinConnection(gameID uuid.UUID, userID int, conn ConnectionWrapper) error
	KickPlayer(gameID uuid.UUID, playerID int) error
}
