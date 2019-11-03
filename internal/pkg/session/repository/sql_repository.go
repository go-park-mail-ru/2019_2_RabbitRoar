package repository

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type sqlSessionRepository struct {
	conn *pgx.Conn
}

func NewSqlSessionRepository(conn *pgx.Conn) session.Repository {
	return &sqlSessionRepository{conn}
}

func (repo sqlSessionRepository) GetUser(sessionID uuid.UUID) (*models.User, error) {
	row := repo.conn.QueryRow(context.Background(), "SELECT id, username, password, email, rating, avatar FROM svoyak.User WHERE id = (SELECT User_id FROM svoyak.Session WHERE UUID = '$1');", sessionID)

	var user models.User
	err := row.Scan(&user.UID, &user.Username, &user.Password, &user.Email, &user.Rating, &user.AvatarUrl)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *sqlSessionRepository) Create(user models.User) (*uuid.UUID, error) {
	newUUID, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	commandTag, err := repo.conn.Exec(context.Background(), "INSERT INTO svoyak.Session VALUES ('$1', $2);", newUUID, user.UID)

	if err != nil {
		return nil, err
	}

	if commandTag.RowsAffected() != 1 {
		return nil, errors.New("Unable to create session: Session already exists")
	}

	return &newUUID, nil
}

func (repo *sqlSessionRepository) Destroy(sessionID uuid.UUID) error {
	commandTag, err := repo.conn.Exec(context.Background(), "DELETE FROM svoyak.Session WHERE UUID = '$1';", sessionID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("Unable to destroy session: No session found")
	}

	return nil
}
