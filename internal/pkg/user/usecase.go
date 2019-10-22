package user

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"

type UseCase interface {
	GetByID(userID int) (*models.User, error)
	GetByName(name string) (*models.User, error)
	Create(user models.User) (*models.User, error)
	UpdatePassword(UID int, password string) error
	IsPasswordCorrect(models.User) (*models.User, bool)
}
