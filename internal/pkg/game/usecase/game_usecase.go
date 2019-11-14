package usecase

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("game_handler")

type gameUseCase struct {
	gameRepo  game.Repository
	packRepo  pack.Repository
	sanitizer *bluemonday.Policy
}

func NewGameUseCase(gameRepo game.Repository, packRepo pack.Repository) game.UseCase {
	return &gameUseCase{
		gameRepo:  gameRepo,
		packRepo:  packRepo,
		sanitizer: bluemonday.UGCPolicy(),
	}
}

func (uc *gameUseCase) GetByID(uuid uuid.UUID) (*models.Game, error) {
	return uc.gameRepo.GetByID(uuid)
}

func (uc *gameUseCase) Create(g models.Game, u models.User) error {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	g.UUID = newUUID
	g.PlayersJoined = 0
	g.Creator = u.ID

	return uc.gameRepo.Create(g)
}

func (uc *gameUseCase) Update(g, gUpdate models.Game) (*models.Game, error) {
	return &g, uc.gameRepo.Update(g)
}

func (uc *gameUseCase) Fetch(page int) (*[]models.Game, error) {
	games, err := uc.gameRepo.Fetch(viper.GetInt("internal.page_size"), page)
	if err != nil {
		return nil, err
	}

	for _, game := range *games {
		gamePack, err := uc.packRepo.GetByID(game.PackID)
		if err != nil {
			return nil, err
		}

		game.PackName = gamePack.Name
	}

	return games, nil
}
