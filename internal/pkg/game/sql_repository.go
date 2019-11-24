package game

import (
	"github.com/google/uuid"
)

type SQLRepository interface {
	JoinPlayer(playerID int, gameID uuid.UUID) error
	KickPlayer(playerID int) (uuid.UUID, error)
	GetGameIDByUserID(userID int) (uuid.UUID, error)
	// GetByID(gameID uuid.UUID) (*models.Game, error)
	// GetPlayers(game models.Game) (*[]models.User, error)
	// FetchOrderedByPlayersJoined(desc bool, pageSize, page int) (*[]models.Game, error)
	// Fetch(pageSize, page int) (*[]models.Game, error)
	// Create(game models.Game) error
	// Update(game models.Game) error
	// Delete(gameID uuid.UUID) error
}
