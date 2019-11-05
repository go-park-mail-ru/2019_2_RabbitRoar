package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
)

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
		"SELECT id, username, password, email, rating, avatar FROM svoyak.User WHERE id = $1;",
		userID,
	)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Rating, &user.AvatarUrl)

	return &user, err
}

func (repo *sqlUserRepository) GetByName(name string) (*models.User, error) {
	row := repo.conn.QueryRow(
		context.Background(),
		"SELECT id, username, password, email, rating, avatar FROM svoyak.User WHERE username = '$1';",
		name,
	)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Rating, &user.AvatarUrl)

	return &user, err
}

func (repo *sqlUserRepository) Create(user models.User) (*models.User, error) {
	if _, err := repo.GetByName(user.Username); err == nil {
		return nil, errors.New("Unable to create user: Username already exists")
	}

	idRow := repo.conn.QueryRow(
		context.Background(),
		"INSERT INTO svoyak.User VALUES (DEFAULT, '$1', '$2', '$3', $4, '$5') RETURNING id;",
		user.Username, user.Password, user.Email, user.Rating, user.AvatarUrl,
	)

	err := idRow.Scan(&user.ID)

	return &user, err
}

func (repo *sqlUserRepository) Update(user models.User) error {
	commandTag, err := repo.conn.Exec(
		context.Background(),
		"UPDATE svoyak.User SET username = '$1', password = '$2', email = '$3', rating = $4, avatar = '$5' WHERE id = $6;",
		user.Username, user.Password, user.Email, user.Rating, user.AvatarUrl, user.ID,
	)

	if commandTag.RowsAffected() != 1 {
		return errors.New("Unable to update user: No user found")
	}

	return err
}
