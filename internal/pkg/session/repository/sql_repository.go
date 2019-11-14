package repository

import (
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session"
	"github.com/google/uuid"
)

type sqlSessionRepository struct {
	db *sql.DB
}

func NewSqlSessionRepository(db *sql.DB) session.Repository {
	return &sqlSessionRepository{
		db: db,
	}
}

func (repo *sqlSessionRepository) Create(user models.User) (*string, error) {
	UUID, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	res, err := repo.db.Exec(
		`
			INSERT INTO "svoyak"."Session" ("UUID", "User_id")
			VALUES ($1::varchar, $2::integer);
		`,
		UUID, user.ID,
	)
	if err != nil {
		return nil, err
	}

	c, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	if c != 1 {
		return nil, errors.New("unable to create session: Session already exists")
	}
	sessionID := UUID.String()

	return &sessionID, err
}

func (repo *sqlSessionRepository) Destroy(sessionID string) error {
	res, err := repo.db.Exec(
		`
			DELETE FROM "svoyak"."Session"
			WHERE "UUID" = $1::varchar;
		`,
		sessionID,
	)
	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("unable to destroy session: No session found")
	}

	return err
}
