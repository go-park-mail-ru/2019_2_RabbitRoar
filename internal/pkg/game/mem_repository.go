package game

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/google/uuid"
)

type MemRepository interface {
	Create(gameID uuid.UUID, host models.User) error
	JoinConnection(gameID uuid.UUID, u models.User, conn ConnectionWrapper) error
}
