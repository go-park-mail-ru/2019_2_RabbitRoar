package game

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"

type Repository interface {
	GetByID(int) (*models.Game, error)
	GetPlayers(pack models.Game) (*[]models.User, error)
	Create(pack models.Game) (*models.Game, error)
	Update(pack models.Game) (*models.Game, error)
	Delete(int) error
}
