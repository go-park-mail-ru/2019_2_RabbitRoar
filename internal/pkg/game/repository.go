package game

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"

type Repository interface {
	GetByID(gameID int) (*models.Game, error)
	GetPlayers(game models.Game) (*[]models.User, error)
	FetchOrderedByPlayersJoined(desc bool, pageSize, page int) (*[]models.Game, error)
	Fetch(pageSize, page int) (*[]models.Game, error)
	Create(pack models.Game) (*models.Game, error)
	Update(pack models.Game) error
	Delete(gameID int) error
}
