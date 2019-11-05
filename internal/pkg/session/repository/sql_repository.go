package repository

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type sqlSessionRepository struct {
	conn *pgxpool.Pool
}

func NewSqlSessionRepository(conn *pgxpool.Pool) session.Repository {
	return &sqlSessionRepository{
		conn: conn,
	}
}

func (repo sqlSessionRepository) GetUser(sessionID uuid.UUID) (*models.User, error) {
	row := repo.conn.QueryRow(
		context.Background(),
		`
			SELECT id, username, password, email, rating, avatar
			FROM "svoyak"."User"
			WHERE "id" = (SELECT "User_id" FROM "svoyak"."Session" WHERE "UUID" = $1::varchar);
		`,
		sessionID,
	)

	var u models.User
	uPassword := make([]byte, 45)
	err := row.Scan(&u.ID, &u.Username, &uPassword, &u.Email, &u.Rating, &u.AvatarUrl)

	u.Password = string(uPassword)

	return &u, err
}

func (repo *sqlSessionRepository) Create(user models.User) (*uuid.UUID, error) {
	newUUID, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	commandTag, err := repo.conn.Exec(
		context.Background(),
		`
			INSERT INTO "svoyak"."Session" ("UUID", "User_id")
			VALUES ($1::varchar, $2::integer);
		`,
		newUUID, user.ID,
	)

	if commandTag.RowsAffected() != 1 {
		return nil, errors.New("unable to create session: Session already exists")
	}

	return &newUUID, err
}

func (repo *sqlSessionRepository) Destroy(sessionID uuid.UUID) error {
	commandTag, err := repo.conn.Exec(
		context.Background(),
		`
			DELETE FROM "svoyak"."Session"
			WHERE "UUID" = $1::varchar;
		`,
		sessionID,
	)

	if commandTag.RowsAffected() != 1 {
		return errors.New("unable to destroy session: No session found")
	}

	return err
}
