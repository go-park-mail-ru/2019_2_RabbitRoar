package game

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type Repository interface {
	Create(g *models.Game, packQuestions interface{}, host *models.User) (*models.Game, error)
	Fetch(pageSize int, page int) (*[]models.Game, error)
	GetGameIDByUserID(userID int) (uuid.UUID, error)

	JoinPlayer(u *models.User, gameID uuid.UUID) (*models.Game, error)
	JoinConnection(gameID uuid.UUID, userID int, conn ConnectionWrapper) error
	KickPlayer(playerID int) error
}
