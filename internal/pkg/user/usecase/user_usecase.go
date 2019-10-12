package usecase

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
)

type userUseCase struct {
	repository user.Repository
}

func NewUserUseCase(userRepo user.Repository) user.UseCase {
	return &userUseCase{
		repository: userRepo,
	}
}

func (uc *userUseCase) GetByID(id int) (*models.User, error) {
	return uc.repository.GetByID(id)
}
