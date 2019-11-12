package game

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type UseCase interface {
	GetByID(uuid uuid.UUID) (*models.Game, error)
	Create(g models.Game) (*models.Game, error)
	Update(g, gUpdate models.Game) (*models.Game, error)
	Fetch(page int) (*[]models.Game, error)
}
