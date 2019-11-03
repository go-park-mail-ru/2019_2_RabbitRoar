package repository

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type sqlGameRepository struct {
	conn *pgx.Conn
}

func NewSqlGameRepository(conn *pgx.Conn) game.Repository {
	return &sqlGameRepository{conn}
}

func (repo sqlGameRepository) GetByID(gameID int) (*models.Game, error) {
	row := repo.conn.QueryRow(context.Background(), "SELECT uuid, name, playersCapacity, playersJoined, state, creator FROM svoyak.Game WHERE uuid = $1;", gameID)

	var game models.Game
	err := row.Scan(&game.UUID, &game.Name, &game.PlayersCapacity, &game.PlayersJoined, &game.State, &game.Creator)

	return &game, err
}

func (repo sqlGameRepository) GetPlayers(game models.Game) (*[]models.User, error) {
	return nil, errors.New("Function not implemented")
}

func (repo sqlGameRepository) FetchOrderedByPlayersJoined(desc bool, pageSize, page int) (*[]models.Game, error) {
	var order string
	if desc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	rows, err := repo.conn.Query(context.Background(), "SELECT uuid, name, playersCapacity, playersJoined, state, creator FROM svoyak.Game ORDER BY playersJoined $1;", order)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []models.Game

	for rows.Next() {
		var game models.Game

		err := rows.Scan(&game.UUID, &game.Name, &game.PlayersCapacity, &game.PlayersJoined, &game.State, &game.Creator)

		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return &games, rows.Err()
}

func (repo sqlGameRepository) Fetch(pageSize, page int) (*[]models.Game, error) {
	rows, err := repo.conn.Query(context.Background(), "SELECT uuid, name, playersCapacity, playersJoined, state, creator FROM svoyak.Game OFFSET $1 LIMIT $2;", (page * pageSize), pageSize)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []models.Game

	for rows.Next() {
		var game models.Game

		err := rows.Scan(&game.UUID, &game.Name, &game.PlayersCapacity, &game.PlayersJoined, &game.State, &game.Creator)

		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return &games, rows.Err()
}

func (repo *sqlGameRepository) Create(game models.Game) (*models.Game, error) {
	newUUID, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	game.UUID = newUUID

	commandTag, err := repo.conn.Exec(context.Background(), "INSERT INTO svoyak.Session VALUES ('$1', '$2', $3, $4, $5, $6);", game.UUID, game.Name, game.PlayersCapacity, game.PlayersJoined, game.State, game.Creator)

	if commandTag.RowsAffected() != 1 {
		return nil, errors.New("Unable to create game: Game already exists")
	}

	return &game, err
}

func (repo *sqlGameRepository) Update(game models.Game) error {
	commandTag, err := repo.conn.Exec(context.Background(), "UPDATE svoyak.Pack SET name = '$1', playerCapacity = $2, playerJoined = $3, state = $4, creator = $5 WHERE uuid = '$6';", game.Name, game.PlayersCapacity, game.PlayersJoined, game.State, game.Creator, game.UUID)

	if commandTag.RowsAffected() != 1 {
		return errors.New("Unable to update game: No game found")
	}

	return err
}

func (repo *sqlGameRepository) Delete(gameID int) error {
	commandTag, err := repo.conn.Exec(context.Background(), "DELETE FROM svoyak.Game WHERE id = $1;", gameID)

	if commandTag.RowsAffected() != 1 {
		return errors.New("Unable to delete game: No game found")
	}

	return err
}
