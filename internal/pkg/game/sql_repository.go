package game

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type SQLRepository interface {
	GetByID(gameID uuid.UUID) (*models.Game, error)
	GetPlayers(game models.Game) (*[]models.User, error)
	GetGameIDByUserID(userID int) (uuid.UUID, error)
	FetchOrderedByPlayersJoined(desc bool, pageSize, page int) (*[]models.Game, error)
	Fetch(pageSize, page int) (*[]models.Game, error)
	JoinPlayer(playerID int, gameID uuid.UUID) error
	KickPlayer(playerID int) (uuid.UUID, error)
	// FetchAllReadyGames() (*[]models.Game, error)
	Create(game models.Game) error
	Update(game models.Game) error
	Delete(gameID uuid.UUID) error
}
