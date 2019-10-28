package session

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type Repository interface {
	GetUser(sessionID uuid.UUID) (*models.User, error)
	Create(user models.User) (*uuid.UUID, error)
	Destroy(sessionID uuid.UUID) error
}
