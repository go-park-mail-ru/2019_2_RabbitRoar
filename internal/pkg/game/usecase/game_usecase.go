package usecase

import (
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("game_handler")

type gameUseCase struct {
	gameRepo  game.Repository
	sanitizer *bluemonday.Policy
}

func NewGameUseCase(gameRepo game.Repository) game.UseCase {
	return &gameUseCase{
		gameRepo:  gameRepo,
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
	g.Pending = true

	return uc.gameRepo.Create(g)
}

func (uc *gameUseCase) Fetch(page int) (*[]models.Game, error) {
	games, err := uc.gameRepo.Fetch(viper.GetInt("internal.page_size"), page)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func (uc *gameUseCase) JoinPlayerToGame(playerID int, gameID uuid.UUID) error {
	game, err := uc.gameRepo.GetByID(gameID)
	if err != nil {
		return err
	}

	if game.PlayersJoined >= game.PlayersCapacity {
		return errors.New("unable to join the room: room is full")
	}

	game.PlayersJoined++

	if game.PlayersJoined == game.PlayersCapacity {
		game.Pending = false
	}

	uc.gameRepo.Update(*game)

	return uc.gameRepo.JoinPlayer(playerID, game.UUID)
}

func (uc *gameUseCase) KickPlayerFromGame(playerID int) error {
	gameID, err := uc.gameRepo.KickPlayer(playerID)
	if err != nil {
		return err
	}

	game, err := uc.gameRepo.GetByID(gameID)
	if err != nil {
		return err
	}

	if game.PlayersJoined <= 0 {
		return errors.New("unable to leave the room: room is empty")
	}

	game.PlayersJoined--
	uc.gameRepo.Update(*game)

	return nil
}
