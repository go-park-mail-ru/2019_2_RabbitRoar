package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type memSessionRepository struct {
	sessions    map[uuid.UUID]models.Session
	lastUserUID int64
}

//func NewMemSessionRepository() session.Repository {
//	return &memSessionRepository{
//		sessions:    map[uuid.UUID]models.Session{},
//		lastUserUID: 0,
//	}
//}

func (repo memSessionRepository) GetUser(sessionID uuid.UUID) (*models.User, error) {
	if s, success := repo.sessions[sessionID]; success {
		user := s.User
		return &user, nil
	}
	return nil, errors.New("session not found")
}

func (repo *memSessionRepository) Create(user models.User) (*uuid.UUID, error) {
	newUUID, err := uuid.NewUUID()

	repo.sessions[newUUID] = models.Session{
		ID: newUUID.String(),
		User: user,
	}

	return &newUUID, err
}

func (repo *memSessionRepository) Destroy(sessionID uuid.UUID) error {
	delete(repo.sessions, sessionID)

	return nil
}
