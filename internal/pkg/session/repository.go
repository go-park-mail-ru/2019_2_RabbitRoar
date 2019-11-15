package session

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
)

type Repository interface {
	Create(user models.User) (*string, error)
	Destroy(sessionID string) error
}
