package session

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
)

type UseCase interface {
	Create(models.User) (*string, error)
	Destroy(string)
}
