package repository

import (
	"database/sql"
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
)

type sqlUserRepository struct {
	db *sql.DB
}

func NewSqlUserRepository(db *sql.DB) user.Repository {
	return &sqlUserRepository{
		db: db,
	}
}

func scanUser(row *sql.Row) (*models.User, error) {
	var u models.User
	password := make([]byte, 36)

	err := row.Scan(
		&u.ID,
		&u.Username,
		&password,
		&u.Email,
		&u.Rating,
		&u.AvatarUrl,
	)

	if err != nil {
		return nil, err
	}

	u.Password = string(password)

	return &u, err
}

func (repo *sqlUserRepository) GetByID(userID int) (*models.User, error) {
	row := repo.db.QueryRow(
		`
			SELECT id, username, password, email, rating, avatar
			FROM "svoyak"."User"
			WHERE id = $1::integer;
		`,
		userID,
	)

	return scanUser(row)
}

func (repo *sqlUserRepository) GetByName(name string) (*models.User, error) {
	row := repo.db.QueryRow(
		`
			SELECT id, username, password, email, rating, avatar
			FROM "svoyak"."User"
			WHERE username = $1::varchar;
		`,
		name,
	)

	return scanUser(row)
}

func (repo *sqlUserRepository) Create(u models.User) (*models.User, error) {
	idRow := repo.db.QueryRow(`
			INSERT INTO "svoyak"."User" (id, username, password, email, rating, avatar)
			VALUES (DEFAULT, $1::varchar, $2::bytea, $3::varchar, $4::integer, $5::varchar)
			RETURNING id;
		`,
		u.Username, []byte(u.Password), u.Email, u.Rating, u.AvatarUrl,
	)

	err := idRow.Scan(&u.ID)

	return &u, err
}

func (repo *sqlUserRepository) Update(user models.User) error {
	res, err := repo.db.Exec(
		`
			UPDATE "svoyak"."User"
			SET username = $1::varchar, password = $2::bytea, email = $3::varchar, rating = $4::integer, avatar = $5::varchar
			WHERE id = $6;
		`,
		user.Username, []byte(user.Password), user.Email, user.Rating, user.AvatarUrl, user.ID,
	)

	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("unable to update user: no user found")
	}

	return nil
}

func (repo *sqlUserRepository) FetchLeaderBoard(page, pageSize int) ([]models.User, error) {
	rows, err := repo.db.Query(
		`
			SELECT id, username, rating, avatar
			FROM "svoyak"."User"
			ORDER BY rating DESC
			OFFSET $1::integer LIMIT $2::integer;
		`,
		page * pageSize, pageSize,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching leaderboard")
	}

	var users = make([]models.User, 0, pageSize)
	for rows.Next() {
		var u models.User
		err :=  rows.Scan(
			&u.ID,
			&u.Username,
			&u.Rating,
			&u.AvatarUrl,
		)

		if err != nil {
			return nil, errors.Wrap(err, "error scanning user")
		}

		users = append(users, u)
	}

	return users, nil
}
