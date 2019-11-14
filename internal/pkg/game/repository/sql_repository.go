package repository

import (
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type sqlGameRepository struct {
	db *sql.DB
}

func NewSqlGameRepository(db *sql.DB) game.Repository {
	return &sqlGameRepository{
		db: db,
	}
}

func (repo sqlGameRepository) GetByID(gameID uuid.UUID) (*models.Game, error) {
	row := repo.db.QueryRow(
		`
			SELECT UUID, name, players_cap, players_joined, creator, Pack_id
			FROM "svoyak"."Game"
			WHERE "UUID" = $1::varchar;
		`,
		gameID,
	)

	var game models.Game
	err := row.Scan(&game.UUID, &game.Name, &game.PlayersCapacity, &game.PlayersJoined, &game.Creator, &game.PackID)

	return &game, err
}

func (repo sqlGameRepository) GetPlayers(game models.Game) (*[]models.User, error) {
	rows, err := repo.db.Query(
		`
			SELECT id, username, password, email, rating, avatar
			FROM "svoyak"."User"
			WHERE "id" = ANY(SELECT User_id
				FROM "svoyak"."GameUser"
				WHERE "Game_UUID" = $1::varchar);
		`,
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

	rows, err := repo.db.Query(
		`
			SELECT UUID, name, players_cap, players_joined, creator, Pack_id
			FROM "svoyak"."Game"
			ORDER BY players_joined $1::text;
		`,
		order,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []models.Game

	for rows.Next() {
		var game models.Game

		err := rows.Scan(&game.UUID, &game.Name, &game.PlayersCapacity, &game.PlayersJoined, &game.Creator, &game.PackID)

		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return &games, rows.Err()
}

func (repo sqlGameRepository) Fetch(pageSize, page int) (*[]models.Game, error) {
	rows, err := repo.db.Query(
		`
			SELECT UUID, name, players_cap, players_joined, creator, Pack_id
			FROM "svoyak"."Game"
			OFFSET $1::integer LIMIT $2::integer;
		`,
		(page * pageSize), pageSize,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []models.Game

	for rows.Next() {
		var game models.Game

		err := rows.Scan(&game.UUID, &game.Name, &game.PlayersCapacity, &game.PlayersJoined, &game.Creator, &game.PackID)

		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return &games, rows.Err()
}

func (repo *sqlGameRepository) Create(game models.Game) error {
	res, err := repo.db.Exec(
		`
			INSERT INTO "svoyak"."Game" (UUID, name, players_cap, players_joined, creator, Pack_id)
			VALUES ($1::varchar, $2::varchar, $3::integer, $4::integer, $5::integer, $6::integer);
		`,
		game.UUID, game.Name, game.PlayersCapacity, game.PlayersJoined, game.Creator, game.PackID,
	)

	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if c != 1 {
		return errors.New("Unable to create game: Game already exists")
	}

	return err
}

func (repo *sqlGameRepository) Update(game models.Game) error {
	res, err := repo.db.Exec(
		`
			UPDATE "svoyak"."Game"
			SET name = $1::varchar, players_cap = $2::integer, players_joined = $3::integer, creator = $4::integer, Pack_id = $5::integer
			WHERE "UUID" = $6::varchar;"
		`,
		game.Name, game.PlayersCapacity, game.PlayersJoined, game.Creator, game.PackID, game.UUID,
	)

	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if c != 1 {
		return errors.New("Unable to update game: No game found")
	}

	return err
}

func (repo *sqlGameRepository) Delete(gameID int) error {
	res, err := repo.db.Exec(
		`
			DELETE FROM "svoyak"."Game"
			WHERE id = $1::integer;
		`,
		gameID,
	)

	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if c != 1 {
		return errors.New("Unable to delete game: No game found")
	}

	return err
}
