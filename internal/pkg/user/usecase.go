package user

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"

type UseCase interface {
	GetByID(int) (*models.User, error)
}
