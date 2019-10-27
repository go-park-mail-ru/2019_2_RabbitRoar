package session

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type Repository interface {
	GetUser(uuid.UUID) (*models.User, error)
	Create(models.User) (*uuid.UUID, error)
	Destroy(sessionID uuid.UUID) error
}
