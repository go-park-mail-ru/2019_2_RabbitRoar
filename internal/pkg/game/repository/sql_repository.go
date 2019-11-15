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

func NewSqlGameRepository(db *sql.DB) game.SQLRepository {
	return &sqlGameRepository{
		db: db,
	}
}

func (repo sqlGameRepository) GetByID(gameID uuid.UUID) (*models.Game, error) {
	row := repo.db.QueryRow(`
			SELECT
				g.UUID,
				g.name,
				g.players_cap,
				g.players_joined,
				g.creator,
				g.pending,
				g.Pack_id,
				p.name
			FROM "svoyak"."Game" g
			INNER JOIN "svoyak"."Pack" p ON g.Pack_id = p.id
			WHERE "g"."UUID" = $1::varchar;
		`,
		gameID,
	)

	var game models.Game
	err := row.Scan(&game.UUID, &game.Name, &game.PlayersCapacity, &game.PlayersJoined, &game.Creator, &game.Pending, &game.PackID, &game.PackName)

	return &game, err
}

func (repo sqlGameRepository) GetPlayers(game models.Game) (*[]models.User, error) {
	rows, err := repo.db.Query(`
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

func (repo sqlGameRepository) GetGameIDByUserID(userID int) (uuid.UUID, error) {
	row := repo.db.QueryRow(`
			SELECT Game_UUID
			FROM "svoyak"."GameUser"
			WHERE "User_id" = $1::integer
		`,
		userID,
	)

	var gameID uuid.UUID

	err := row.Scan(&gameID)

	return gameID, err
}

func (repo sqlGameRepository) FetchOrderedByPlayersJoined(desc bool, pageSize, page int) (*[]models.Game, error) {
	var order string
	if desc {
		order = "DESC"
	} else {
		order = "ASC"
	}

	rows, err := repo.db.Query(`
			SELECT
				g.UUID,
				g.name,
				g.players_cap,
				g.players_joined,
				g.creator,
				g.pending,
				g.Pack_id,
				p.name
			FROM "svoyak"."Game" g
			INNER JOIN "svoyak"."Pack" p ON g.Pack_id = p.id
			WHERE "g"."pending" = TRUE
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

		err := rows.Scan(&game.UUID, &game.Name, &game.PlayersCapacity, &game.PlayersJoined, &game.Creator, &game.Pending, &game.PackID, &game.PackName)

		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return &games, rows.Err()
}

func (repo sqlGameRepository) Fetch(pageSize, page int) (*[]models.Game, error) {
	rows, err := repo.db.Query(`
			SELECT 
				g.UUID,
				g.name,
				g.players_cap,
				g.players_joined,
				g.creator,
				g.pending,
				g.Pack_id,
				p.name
			FROM
				"svoyak"."Game" g
			INNER JOIN "svoyak"."Pack" p ON g.Pack_id = p.id
			WHERE "g"."pending" = TRUE
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

		err := rows.Scan(&game.UUID, &game.Name, &game.PlayersCapacity, &game.PlayersJoined, &game.Creator, &game.Pending, &game.PackID, &game.PackName)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return &games, rows.Err()
}

func (repo *sqlGameRepository) JoinPlayer(userID int, gameID uuid.UUID) error {
	res, err := repo.db.Exec(`
			INSERT INTO "svoyak"."GameUser" (User_id, Game_UUID)
			VALUES ($1::integer, $2::varchar);
		`,
		userID, gameID,
	)

	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if c != 1 {
		return errors.New("Unable to join game: User already joined")
	}

	return err
}

func (repo *sqlGameRepository) KickPlayer(playerID int) (uuid.UUID, error) {
	row := repo.db.QueryRow(`
			DELETE FROM "svoyak"."GameUser"
			WHERE "User_id" = $1::integer
			RETURNING Game_UUID;
		`,
		playerID,
	)

	var gameID uuid.UUID

	err := row.Scan(&gameID)

	return gameID, err
}

func (repo *sqlGameRepository) Create(game models.Game) error {
	res, err := repo.db.Exec(`
			INSERT INTO "svoyak"."Game" (UUID, name, players_cap, players_joined, creator, pending, Pack_id)
			VALUES ($1::varchar, $2::varchar, $3::integer, $4::integer, $5::integer, $6::boolean, $7::integer);
		`,
		game.UUID, game.Name, game.PlayersCapacity, game.PlayersJoined, game.Creator, game.Pending, game.PackID,
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
	res, err := repo.db.Exec(`
			UPDATE "svoyak"."Game"
			SET name = $1::varchar, players_cap = $2::integer, players_joined = $3::integer, creator = $4::integer, pending = $5::boolean, Pack_id = $6::integer
			WHERE "UUID" = $7::varchar;"
		`,
		game.Name, game.PlayersCapacity, game.PlayersJoined, game.Creator, game.Pending, game.PackID, game.UUID,
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

func (repo *sqlGameRepository) Delete(gameID uuid.UUID) error {
	res, err := repo.db.Exec(`
			DELETE FROM "svoyak"."Game"
			WHERE UUID = $1::varchar;
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
