package usecase

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("game_handler")

type gameUseCase struct {
	repository game.Repository
	sanitizer  *bluemonday.Policy
}

func NewGameUseCase(gameRepo game.Repository) game.UseCase {
	return &gameUseCase{
		repository: gameRepo,
		sanitizer:  bluemonday.UGCPolicy(),
	}
}

func (uc *gameUseCase) GetByID(uuid uuid.UUID) (*models.Game, error) {
	return uc.repository.GetByID(uuid)
}

func (uc *gameUseCase) Create(g models.Game) (*models.Game, error) {
	return uc.repository.Create(g)
}

func (uc *gameUseCase) Update(g, gUpdate models.Game) (*models.Game, error) {
	return &g, uc.repository.Update(g)
}
