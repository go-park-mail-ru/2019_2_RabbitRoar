package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/google/uuid"
)

type playerConnections []game.PlayerConnection

type memGameRepository struct {
	games map[uuid.UUID]playerConnections
}

func NewMemGameRepository() game.MemRepository {
	return &memGameRepository{}
}

func (repo *memGameRepository) JoinConnection(gameID uuid.UUID, conn game.PlayerConnection) error {
	if _, exists := repo.games[gameID]; !exists {
		return errors.New("no game found to join connection")
	}

	repo.games[gameID] = append(repo.games[gameID], conn)

	return nil
}
