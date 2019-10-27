package user

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"

type Repository interface {
	GetByID(userID int) (*models.User, error)
	GetByName(name string) (*models.User, error)
	Create(user models.User) (*models.User, error)
	Update(user models.User) error
}
