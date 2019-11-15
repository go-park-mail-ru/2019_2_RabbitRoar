package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/google/uuid"
)

type gameObj struct {
	hostID       int
	hostConn     game.Connection
	playersConns []game.Connection
}

type memGameRepository struct {
	games map[uuid.UUID]*gameObj
}

func NewMemGameRepository() game.MemRepository {
	return &memGameRepository{
		games: make(map[uuid.UUID]*gameObj),
	}
}

func (repo *memGameRepository) Create(gameID uuid.UUID, hostID int) error {
	if _, exists := repo.games[gameID]; exists {
		return errors.New("game already exists")
	}

	repo.games[gameID] = &gameObj{
		hostID: hostID,
	}

	return nil
}

func (repo *memGameRepository) JoinConnection(gameID uuid.UUID, conn game.Connection) error {
	if _, exists := repo.games[gameID]; !exists {
		return errors.New("no game found to join connection")
	}

	if conn.GetUserID() == repo.games[gameID].hostID {
		repo.games[gameID].hostConn = conn
	} else {
		repo.games[gameID].playersConns = append(repo.games[gameID].playersConns, conn)
	}

	return nil
}
