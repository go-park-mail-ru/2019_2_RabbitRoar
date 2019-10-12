package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type memUserRepository struct {
	sessions    map[uuid.UUID]models.Session
	lastUserUID int64
}

func (repo memUserRepository) SessionGetUser(sessionId uuid.UUID) (*models.User, error) {
	if session, success := repo.sessions[sessionId]; success {
		var user models.User = session.User
		return &user, nil
	}
	return nil, errors.New("Session not found")
}

func (repo *memUserRepository) SessionCreate(user models.User) (uuid.UUID, error) {
	newUUID, err := uuid.NewUUID()

	repo.sessions[newUUID] = models.Session{
		Uuid: newUUID,
		User: user,
	}

	return newUUID, err
}

func (repo *memUserRepository) SessionDestroy(sessionId uuid.UUID) {
	delete(repo.sessions, sessionId)
}
