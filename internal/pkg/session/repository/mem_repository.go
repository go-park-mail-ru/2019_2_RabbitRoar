package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type memSessionRepository struct {
	sessions    map[uuid.UUID]models.Session
	lastUserUID int64
}

func NewMemSessionRepository() session.Repository {
	return &memSessionRepository{
		sessions:    map[uuid.UUID]models.Session{},
		lastUserUID: 0,
	}
}

func (repo memSessionRepository) GetUser(sessionId uuid.UUID) (*models.User, error) {
	if s, success := repo.sessions[sessionId]; success {
		user := s.User
		return &user, nil
	}
	return nil, errors.New("session not found")
}

func (repo *memSessionRepository) Create(user models.User) (*uuid.UUID, error) {
	newUUID, err := uuid.NewUUID()

	repo.sessions[newUUID] = models.Session{
		Uuid: newUUID,
		User: user,
	}

	return &newUUID, err
}

func (repo *memSessionRepository) Destroy(sessionId uuid.UUID) {
	delete(repo.sessions, sessionId)
}
