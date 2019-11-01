package usecase

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/auth"
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

func (uc *userUseCase) Create(u models.User) (*models.User, error) {
	if ok, err := govalidator.ValidateStruct(u); !ok {
		return nil, err
	}

	if err := uc.prepare(&u); err != nil {
		return nil, err
	}

	return uc.repository.Create(u)
}

func (uc *userUseCase) prepare(u *models.User) error {
	ok, err := govalidator.ValidateStruct(u)
	if !ok {
		return err
	}

	if err := uc.preparePassword(u); err != nil {
		return err
	}

	if err := uc.prepareUsername(u); err != nil {
		return err
	}

	return nil
}

func (uc *userUseCase) preparePassword(u *models.User) error {
	u.Password = auth.HashPassword(u.Password)
	return nil
}

func (uc *userUseCase) prepareUsername(u *models.User) error {
	//TODO: validate username here?
	//TODO: check username existance or just return err in repository save?
	return nil
}

func (uc *userUseCase) Update(u, uUpdate models.User) error {
	if uUpdate.Password != "" {
		u.Password = uUpdate.Password
	}

	if uUpdate.Username != "" {
		u.Password = uUpdate.Password
	}

	if err := uc.prepare(&u); err != nil {
		return err
	}

	return uc.repository.Update(u)
}

func (uc *userUseCase) GetByID(id int) (*models.User, error) {
	return uc.repository.GetByID(id)
}

func (uc *userUseCase) GetByName(name string) (*models.User, error) {
	return uc.repository.GetByName(name)
}

func (uc *userUseCase) IsPasswordCorrect(u models.User) (*models.User, bool) {
	correctUser, err := uc.repository.GetByName(u.Username)
	if err != nil {
		return nil, false
	}

	if !auth.CheckPassword(u.Password, correctUser.Password) {
		return nil, false
	}

	return correctUser, true
}
