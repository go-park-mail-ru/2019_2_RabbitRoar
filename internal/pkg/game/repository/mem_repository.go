package repository

import (
	"errors"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
	"github.com/google/uuid"
)

type memGameRepository struct {
	userRepo user.Repository
	games map[uuid.UUID]*game.Game
	userGame map[int]uuid.UUID
	gameKillerChan chan uuid.UUID
}

func NewMemGameRepository(userRepo user.Repository) game.Repository {
	repo := &memGameRepository{
		userRepo: userRepo,
		games: make(map[uuid.UUID]*game.Game),
		userGame: make(map[int]uuid.UUID),
		gameKillerChan: make(chan uuid.UUID, 10),
	}

	go repo.runGameKiller()

	return repo
}

func (repo *memGameRepository) Create(g *models.Game, packQuestions interface{}, host *models.User) (*models.Game, error) {
	if _, exists := repo.userGame[host.ID]; exists {
		return nil, errors.New("user is already playing")
	}

	if _, exists := repo.games[g.UUID]; exists {
		return nil, errors.New("game already exists")
	}

	repo.games[g.UUID] = &game.Game{
		Players:   []*game.Player{},
		Model:     *g,
		Questions: game.NewQuestionTable(packQuestions),
		EvChan:    make(chan game.EventWrapper, 50),
		Started:   false,
		UserRepo:  repo.userRepo,
	}

	defer func(){
		go repo.games[g.UUID].Run(repo.gameKillerChan)
	}()

	return repo.JoinPlayer(host, g.UUID)
}

func (repo *memGameRepository) Fetch(pageSize int, page int) (*[]models.Game, error) {
	gSlice := make([]models.Game, 0)

	firstEntry := pageSize * page
	iter := 0

	for _, g := range repo.games {
		if g.Started {
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

func (repo *memGameRepository) GetGameIDByUserID(userID int) (uuid.UUID, error) {
	if _, exists := repo.userGame[userID]; !exists {
		return uuid.UUID{}, errors.New("no playing user found")
	}

	return repo.userGame[userID], nil
}

func (repo *memGameRepository) JoinPlayer(u *models.User, gameID uuid.UUID) (*models.Game, error) {
	if _, exists := repo.games[gameID]; !exists {
		return nil, errors.New("no game found")
	}

	if _, exists := repo.userGame[u.ID]; exists {
		return nil, errors.New("user is already in a game")
	}

	if repo.games[gameID].Model.PlayersJoined >= repo.games[gameID].Model.PlayersCapacity {
		return nil, errors.New("unable to join the game: game is full")
	}

	repo.userGame[u.ID] = gameID

	repo.games[gameID].Players = append(
		repo.games[gameID].Players,
		&game.Player{
			Info: game.PlayerInfo{
				ID:       u.ID,
				Username: u.Username,
				Avatar:   u.AvatarUrl,
				Score:    0,
			},
		},
	)

	if repo.games[gameID].Host == nil {
		repo.games[gameID].Host = repo.games[gameID].Players[0]
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

func (repo *memGameRepository) KickPlayer(playerID int) error {
	gameID, exists := repo.userGame[playerID]
	if !exists {
		return errors.New("player is not in game")
	}

	delete(repo.userGame, playerID)

	if _, exists := repo.games[gameID]; !exists {
		return errors.New("no game found to leave")
	}

	for i, p := range repo.games[gameID].Players {
		if p.Info.ID == playerID {
			player := repo.games[gameID].Players[i]

			repo.games[gameID].UpdateUserRating(player.Info)

			repo.games[gameID].Players = append(repo.games[gameID].Players[:i], repo.games[gameID].Players[i+1:]...)

			repo.games[gameID].Model.PlayersJoined--

			if player.Conn != nil {
				player.Conn.Stop()
			}

			repo.games[gameID].EvChan <- game.EventWrapper{
				SenderID: p.Info.ID,
				Event:    &game.Event{
					Type:    game.PlayerLeft,
					Payload: nil,
				},
			}

			return nil
		}
	}

	return errors.New("no player found to kick from game")
}

func (repo *memGameRepository) runGameKiller() {
	for {
		gameID := <-repo.gameKillerChan

		for _, p := range repo.games[gameID].Players {
			repo.KickPlayer(p.Info.ID)
		}

		delete(repo.games, gameID)
	}
}