package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type player struct {
	info game.PlayerInfo
	conn game.Connection
}

type gameObj struct {
	host    player
	players []player
}

type memGameRepository struct {
	games map[uuid.UUID]*gameObj
}

func NewMemGameRepository() game.MemRepository {
	return &memGameRepository{
		games: make(map[uuid.UUID]*gameObj),
	}
}

func (repo *memGameRepository) Create(gameID uuid.UUID, host models.User) error {
	if _, exists := repo.games[gameID]; exists {
		return errors.New("game already exists")
	}

	repo.games[gameID] = &gameObj{
		host: player{
			info: game.PlayerInfo{
				ID:       host.ID,
				Username: host.Username,
				Avatar:   host.AvatarUrl,
				Score:    0,
			},
		},
	}

	return nil
}

func (repo *memGameRepository) JoinConnection(gameID uuid.UUID, u models.User, conn game.Connection) error {
	if _, exists := repo.games[gameID]; !exists {
		return errors.New("no game found to join connection")
	}

	if u.ID == repo.games[gameID].host.info.ID {
		repo.games[gameID].host.conn = conn
	} else {
		repo.games[gameID].players = append(
			repo.games[gameID].players,
			player{
				info: game.PlayerInfo{
					ID:       u.ID,
					Username: u.Username,
					Avatar:   u.AvatarUrl,
					Score:    0,
				},
				conn: conn,
			},
		)
	}

	return nil
}
