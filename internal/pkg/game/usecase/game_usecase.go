package usecase

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game/connection"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

type gameUseCase struct {
	gameSQLRepo game.SQLRepository
	gameMemRepo game.MemRepository
	packRepo    pack.Repository
	sanitizer   *bluemonday.Policy
}

func NewGameUseCase(gameSQLRepo game.SQLRepository, gameMemRepo game.MemRepository, packRepo pack.Repository) game.UseCase {
	return &gameUseCase{
		gameSQLRepo: gameSQLRepo,
		gameMemRepo: gameMemRepo,
		packRepo:    packRepo,
		sanitizer:   bluemonday.UGCPolicy(),
	}
}

func (uc *gameUseCase) Create(g *models.Game, u models.User) error {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	g.UUID = newUUID
	g.PlayersJoined = 0

	p, err := uc.packRepo.GetByID(g.PackID)
	if err != nil {
		return err
	}
	g.PackName = p.Name

	return uc.gameMemRepo.Create(g, u)
}

func (uc *gameUseCase) Fetch(page int) (*[]models.Game, error) {
	return uc.gameMemRepo.Fetch(viper.GetInt("internal.page_size"), page)
}

func (uc *gameUseCase) JoinPlayerToGame(u models.User, gameID uuid.UUID) (*models.Game, error) {
	err := uc.gameSQLRepo.JoinPlayer(u.ID, gameID)
	if err != nil {
		return nil, err
	}

	return uc.gameMemRepo.JoinPlayer(u, gameID)
}

func (uc *gameUseCase) KickPlayerFromGame(playerID int) error {
	gameID, err := uc.gameSQLRepo.KickPlayer(playerID)
	if err != nil {
		return err
	}

	return uc.gameMemRepo.KickPlayer(gameID, playerID)
}

func (uc *gameUseCase) GetGameIDByUserID(userID int) (uuid.UUID, error) {
	return uc.gameSQLRepo.GetGameIDByUserID(userID)
}

func (uc *gameUseCase) NewConnectionWrapper(ws *websocket.Conn) game.ConnectionWrapper {
	sendChan := make(chan game.Event, 5)
	stop := make(chan bool, 2)

	return connection.NewConnectionWrapper(ws, sendChan, stop)
}

func (uc *gameUseCase) JoinConnectionToGame(gameID uuid.UUID, userID int, conn game.ConnectionWrapper) error {
	return uc.gameMemRepo.JoinConnection(gameID, userID, conn)
}
