package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game/state"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type memGameRepository struct {
	games map[uuid.UUID]*game.Game
}

func NewMemGameRepository() game.MemRepository {
	return &memGameRepository{
		games: make(map[uuid.UUID]*game.Game),
	}
}

func (repo *memGameRepository) Create(gameID uuid.UUID, host models.User) error {
	if _, exists := repo.games[gameID]; exists {
		return errors.New("game already exists")
	}

	repo.games[gameID] = &game.Game{
		Host: game.Player{
			Info: game.PlayerInfo{
				ID:       host.ID,
				Username: host.Username,
				Avatar:   host.AvatarUrl,
				Score:    0,
			},
		},
		State: &state.PendPlayers{},
	}

	go repo.games[gameID].Run()

	return nil
}

func (repo *memGameRepository) JoinConnection(gameID uuid.UUID, u models.User, conn game.ConnectionWrapper) error {
	if _, exists := repo.games[gameID]; !exists {
		return errors.New("no game found to join connection")
	}

	if u.ID == repo.games[gameID].Host.Info.ID {
		repo.games[gameID].Host.Conn = conn
	} else {
		repo.games[gameID].Players = append(
			repo.games[gameID].Players,
			game.Player{
				Info: game.PlayerInfo{
					ID:       u.ID,
					Username: u.Username,
					Avatar:   u.AvatarUrl,
					Score:    0,
				},
				Conn: conn,
			},
		)
	}

	return nil
}
