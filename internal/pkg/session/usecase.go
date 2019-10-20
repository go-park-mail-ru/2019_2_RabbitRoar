package session

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type UseCase interface {
	GetUserByUUID(uuid.UUID) (*models.User, error)
	GetUserByStringUUID(string) (*models.User, error)
	Create(models.User) (*uuid.UUID, error)
	Destroy(uuid.UUID)
}
