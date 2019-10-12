package user

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"

type Repository interface {
	GetByID(int) (*models.User, error)
	GetByName(string) (*models.User, error)
	Create(string, string, string) (*models.User, error)
	Update(models.User) error
}
