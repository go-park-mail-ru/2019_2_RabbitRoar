package usecase

import (
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game/connection"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("game_handler")

type gameUseCase struct {
	sqlRepo   game.SQLRepository
	memRepo   game.MemRepository
	sanitizer *bluemonday.Policy
}

func NewGameUseCase(sqlRepo game.SQLRepository, memRepo game.MemRepository) game.UseCase {
	return &gameUseCase{
		sqlRepo:   sqlRepo,
		memRepo:   memRepo,
		sanitizer: bluemonday.UGCPolicy(),
	}
}

func (uc *gameUseCase) GetByID(uuid uuid.UUID) (*models.Game, error) {
	return uc.sqlRepo.GetByID(uuid)
}

func (uc *gameUseCase) GetGameIDByUserID(userID int) (uuid.UUID, error) {
	return uc.sqlRepo.GetGameIDByUserID(userID)
}

func (uc *gameUseCase) SQLCreate(g models.Game, u models.User) error {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	g.UUID = newUUID
	g.PlayersJoined = 0
	g.Creator = u.ID
	g.Pending = true

	return uc.sqlRepo.Create(g)
}

func (uc *gameUseCase) Fetch(page int) (*[]models.Game, error) {
	games, err := uc.sqlRepo.Fetch(viper.GetInt("internal.page_size"), page)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func (uc *gameUseCase) JoinPlayerToGame(playerID int, gameID uuid.UUID) error {
	game, err := uc.sqlRepo.GetByID(gameID)
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

	uc.sqlRepo.Update(*game)

	return uc.sqlRepo.JoinPlayer(playerID, game.UUID)
}

func (uc *gameUseCase) KickPlayerFromGame(playerID int) error {
	gameID, err := uc.sqlRepo.KickPlayer(playerID)
	if err != nil {
		return err
	}

	game, err := uc.sqlRepo.GetByID(gameID)
	if err != nil {
		return err
	}

	if game.PlayersJoined <= 0 {
		return errors.New("unable to leave the room: room is empty")
	}

	game.PlayersJoined--

	if game.PlayersJoined <= 0 {
		uc.sqlRepo.Delete(game.UUID)
	} else {
		uc.sqlRepo.Update(*game)
	}

	return nil
}

func (uc *gameUseCase) NewConnection(userID int) game.Connection {
	sendChan := make(chan game.Event, 5)
	receiveChan := make(chan game.Event, 5)
	stopSend := make(chan bool)
	stopReceive := make(chan bool)

	return connection.NewConnection(userID, sendChan, receiveChan, stopSend, stopReceive)
}

func (uc *gameUseCase) MemCreate(g models.Game, u models.User) error {
	return uc.memRepo.Create(g.UUID, u.ID)
}

func (uc *gameUseCase) JoinConnectionToGame(gameID uuid.UUID, conn game.Connection) error {
	return uc.memRepo.JoinConnection(gameID, conn)
}
