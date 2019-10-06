package repository

import (
	"errors"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/entity"
)

var ErrNotFound = errors.New("no entity found")
var ErrConflict = errors.New("new entity conflicts with existing")

type UserRepository interface {
	UserGetByName(string) (error, entity.User)
	UserCreate(string, string) (error, entity.User)
	UserUpdate(entity.User) error
}

type SessionRepository interface {
	SessionGetUser(uuid.UUID) (entity.User, error)
	SessionCreate(entity.User) (uuid.UUID, error)
	SessionDestroy(sessionId uuid.UUID)
}

type Repository interface {
	UserRepository
	SessionRepository
}
