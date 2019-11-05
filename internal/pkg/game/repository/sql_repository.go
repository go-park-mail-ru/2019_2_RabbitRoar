package repository

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type sqlGameRepository struct {
	conn *pgxpool.Pool
}

func NewSqlGameRepository(conn *pgxpool.Pool) game.Repository {
	return &sqlGameRepository{
		conn: conn,
	}
}

func (repo sqlGameRepository) GetByID(gameID uuid.UUID) (*models.Game, error) {
	row := repo.conn.QueryRow(
		context.Background(),
		"SELECT UUID, name, players_cap, players_joined, state, creator"+
			"FROM svoyak.\"Game\""+
			"WHERE UUID = $1::varchar;",
		gameID,
	)

	var game models.Game
	err := row.Scan(&game.UUID, &game.Name, &game.PlayersCapacity, &game.PlayersJoined, &game.State, &game.Creator)

	return &game, err
}

func (repo sqlGameRepository) GetPlayers(game models.Game) (*[]models.User, error) {
	rows, err := repo.conn.Query(
		context.Background(),
		"SELECT id, username, password, email, rating, avatar"+
			"FROM svoyak.\"User\""+
			"WHERE id = ANY(SELECT User_id"+
			"FROM svoyak.\"GameUser\""+
			"WHERE Game_UUID = $1::varchar);",
		game.UUID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Rating, &user.AvatarUrl)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, rows.Err()
}

func (repo sqlGameRepository) FetchOrderedByPlayersJoined(desc bool, pageSize, page int) (*[]models.Game, error) {
	var order string
	if desc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	rows, err := repo.conn.Query(
		context.Background(),
		"SELECT UUID, name, players_cap, players_joined, state, creator"+
			"FROM svoyak.\"Game\""+
			"ORDER BY players_joined $1::text;",
		order,
	)

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
	rows, err := repo.conn.Query(
		context.Background(),
		"SELECT UUID, name, players_cap, players_joined, state, creator"+
			"FROM svoyak.\"Game\""+
			"OFFSET $1::integer LIMIT $2::integer;",
		(page * pageSize), pageSize,
	)

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

	commandTag, err := repo.conn.Exec(
		context.Background(),
		"INSERT INTO svoyak.\"Session\" (UUID, name, players_cap, players_joined, state, creator)"+
			"VALUES ($1::varchar, $2::varchar, $3::integer, $4::integer, $5::integer, $6::integer);",
		game.UUID, game.Name, game.PlayersCapacity, game.PlayersJoined, game.State, game.Creator,
	)

	if commandTag.RowsAffected() != 1 {
		return nil, errors.New("Unable to create game: Game already exists")
	}

	return &game, err
}

func (repo *sqlGameRepository) Update(game models.Game) error {
	commandTag, err := repo.conn.Exec(
		context.Background(),
		"UPDATE svoyak.Pack"+
			"SET name = $1::varchar, player_cap = $2::integer, player_joined = $3::integer, state = $4::integer, creator = $5::integer"+
			"WHERE uuid = $6::varchar;",
		game.Name, game.PlayersCapacity, game.PlayersJoined, game.State, game.Creator, game.UUID,
	)

	if commandTag.RowsAffected() != 1 {
		return errors.New("Unable to update game: No game found")
	}

	return err
}

func (repo *sqlGameRepository) Delete(gameID int) error {
	commandTag, err := repo.conn.Exec(
		context.Background(),
		"DELETE FROM svoyak.Game"+
			"WHERE id = $1::integer;",
		gameID,
	)

	if commandTag.RowsAffected() != 1 {
		return errors.New("Unable to delete game: No game found")
	}

	return err
}
