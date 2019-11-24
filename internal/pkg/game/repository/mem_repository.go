package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
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

func (repo *memGameRepository) Create(g *models.Game, host models.User) error {
	if _, exists := repo.games[g.UUID]; exists {
		return errors.New("game already exists")
	}

	repo.games[g.UUID] = &game.Game{
		Players: []game.Player{},
		Model:   *g,
		EvChan:  make(chan game.EventWrapper, 50),
	}

	go repo.games[g.UUID].Run()

	return nil
}

func (repo *memGameRepository) Fetch(pageSize int, page int) (*[]models.Game, error) {
	gSlice := make([]models.Game, 0)

	firstEntry := pageSize * page
	iter := 0

	for _, g := range repo.games {
		if g.State.GetType() != game.Pending {
			continue
		}

		if iter < firstEntry {
			iter++
			continue
		}

		gSlice = append(gSlice, g.Model)

		if len(gSlice) == pageSize {
			break
		}
	}

	return &gSlice, nil
}

func (repo *memGameRepository) JoinPlayer(u models.User, gameID uuid.UUID) (*models.Game, error) {
	if _, exists := repo.games[gameID]; !exists {
		return nil, errors.New("no game found")
	}

	if repo.games[gameID].Model.PlayersJoined >= repo.games[gameID].Model.PlayersCapacity {
		return nil, errors.New("unable to join the game: game is full")
	}

	repo.games[gameID].Players = append(
		repo.games[gameID].Players,
		game.Player{
			Info: game.PlayerInfo{
				ID:       u.ID,
				Username: u.Username,
				Avatar:   u.AvatarUrl,
				Score:    0,
			},
		},
	)

	if repo.games[gameID].Host == nil {
		repo.games[gameID].Host = &repo.games[gameID].Players[0]
	}

	repo.games[gameID].Model.PlayersJoined++

	return &repo.games[gameID].Model, nil
}

func (repo *memGameRepository) JoinConnection(gameID uuid.UUID, userID int, conn game.ConnectionWrapper) error {
	if _, exists := repo.games[gameID]; !exists {
		return errors.New("no game found to join player")
	}

	for i, p := range repo.games[gameID].Players {
		if p.Info.ID == userID {
			repo.games[gameID].Players[i].Conn = conn
			repo.games[gameID].Players[i].Conn.SetReceiveChan(repo.games[gameID].EvChan)
			return nil
		}
	}

	return errors.New("no player found to join connection")
}

func (repo *memGameRepository) KickPlayer(gameID uuid.UUID, playerID int) error {
	if _, exists := repo.games[gameID]; !exists {
		return errors.New("no game found to leave")
	}

	for i, p := range repo.games[gameID].Players {
		if p.Info.ID == playerID {
			if repo.games[gameID].Players[i].Conn != nil {
				repo.games[gameID].Players[i].Conn.Stop()
			}
			repo.games[gameID].Players = append(repo.games[gameID].Players[:i], repo.games[gameID].Players[i+1:]...)
			repo.games[gameID].Model.PlayersJoined--
			return nil
		}
	}

	return errors.New("no player found to kick from game")
}
