package repository

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
	"github.com/jackc/pgx/v4"
)

type sqlUserRepository struct {
	conn *pgx.Conn
}

func NewSqlUserRepository(conn *pgx.Conn) user.Repository {
	return &sqlUserRepository{conn}
}

func interpretUser(rows *pgx.Rows) (models.User, error) {
	var user models.User

	for (*rows).Next() {
		err := (*rows).Scan(&user.UID, &user.Username, &user.Password, &user.Email, &user.Rating, &user.AvatarUrl)
		if err != nil {
			return user, err
		}
	}

	return user, (*rows).Err()
}

func (repo *sqlUserRepository) GetByID(userID int) (*models.User, error) {
	rows, err := repo.conn.Query(context.Background(), "SELECT id, username, password, email, rating, avatar FROM svoyak.User WHERE id = $1", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user, err := interpretUser(&rows)

	return &user, err
}

func (repo *sqlUserRepository) GetByName(name string) (*models.User, error) {
	rows, err := repo.conn.Query(context.Background(), "SELECT id, username, password, email, rating, avatar FROM svoyak.User WHERE username = $1", name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user, err := interpretUser(&rows)

	return &user, err
}

func (repo *sqlUserRepository) Create(user models.User) (*models.User, error) {
	if _, err := repo.GetByName(user.Username); err == nil {
		return nil, errors.New("Unable to create user: Username already exists")
	}

	_, err := repo.conn.Exec(context.Background(), "INSERT INTO svoyak.User VALUES (DEFAULT, '$1', '$2', '$3', $4, '$5')", user.Username, user.Password, user.Email, user.Rating, user.AvatarUrl)

	if err != nil {
		return nil, err
	}

	return repo.GetByName(user.Username)
}

func (repo *sqlUserRepository) Update(user models.User) error {
	_, err := repo.conn.Exec(context.Background(), "UPDATE svoyak.User SET username = $1, password = $2, email = $3, rating = $4, avatar = $5 WHERE id = $6", user.Username, user.Password, user.Email, user.Rating, user.AvatarUrl, user.UID)

	return err
}
