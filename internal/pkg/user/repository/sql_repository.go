package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
)

//TODO: REMOVE USER PARSING COPYPASTE

type sqlUserRepository struct {
	conn *pgxpool.Pool
}

func NewSqlUserRepository(conn *pgxpool.Pool) user.Repository {
	return &sqlUserRepository{
		conn: conn,
	}
}

func (repo *sqlUserRepository) GetByID(userID int) (*models.User, error) {
	row := repo.conn.QueryRow(
		context.Background(),
		`
			SELECT id, username, password, email, rating, avatar
			FROM "svoyak"."User"
			WHERE id = $1::integer;
		`,
		userID,
	)

	var u models.User
	uPassword := make([]byte, 45)
	err := row.Scan(&u.ID, &u.Username, &uPassword, &u.Email, &u.Rating, &u.AvatarUrl)

	u.Password = string(uPassword)

	return &u, err
}

func (repo *sqlUserRepository) GetByName(name string) (*models.User, error) {
	row := repo.conn.QueryRow(
		context.Background(),
		`
			SELECT id, username, password, email, rating, avatar
			FROM "svoyak"."User"
			WHERE username = $1::varchar;
		`,
		name,
	)

	var u models.User
	uPassword := make([]byte, 45)
	err := row.Scan(&u.ID, &u.Username, &uPassword, &u.Email, &u.Rating, &u.AvatarUrl)

	u.Password = string(uPassword)

	return &u, err
}

func (repo *sqlUserRepository) Create(u models.User) (*models.User, error) {
	if _, err := repo.GetByName(u.Username); err == nil {
		return nil, errors.New("Unable to create user: Username already exists")
	}

	idRow := repo.conn.QueryRow(
		context.Background(),
		`
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
	commandTag, err := repo.conn.Exec(
		context.Background(),
		`
			UPDATE "svoyak"."User"
			SET username = $1::varchar, password = $2::bytea, email = $3::varchar, rating = $4::integer, avatar = $5::varchar
			WHERE id = $6;
		`,
		user.Username, []byte(user.Password), user.Email, user.Rating, user.AvatarUrl, user.ID,
	)

	if commandTag.RowsAffected() != 1 {
		return errors.New("Unable to update user: No user found")
	}

	return err
}
