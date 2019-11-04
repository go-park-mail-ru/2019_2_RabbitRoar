package user

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"mime/multipart"
)

type UseCase interface {
	GetByID(userID int) (*models.User, error)
	GetByName(name string) (*models.User, error)
	Create(user models.User) (*models.User, error)
	Update(u, uUpdate models.User) error
	UpdateAvatar(u models.User, file *multipart.FileHeader) (*models.User, error)
	IsPasswordCorrect(models.User) (*models.User, bool)
	Sanitize(models.User) models.User
}
