package game

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"

type Repository interface {
	GetByID(ID int) (*models.Game, error)
	GetPlayers(pack models.Game) (*[]models.User, error)
	FetchByPlayersJoined(asc bool, pageSize, page int) (*[]models.Game, error)
	Fetch(pageSize, page int) (*[]models.Game, error)
	Create(pack *models.Game) error
	Update(pack *models.Game) error
	Delete(ID int) error
}
