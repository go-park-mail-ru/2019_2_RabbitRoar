package usecase

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"
)

//TODO: make password saved and validated in hash!
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

	userCreated, err := uc.repository.Create(u)

	if err != nil {
		return nil, err
	}

	return userCreated, nil
}

func (uc *userUseCase) UpdatePassword(UID int, password string) error {
	u, err := uc.repository.GetByID(UID)
	if err != nil {
		return err
	}
	//TODO: validate password here
	u.Password = password

	return uc.repository.Update(*u)
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

	if correctUser.Password != u.Password {
		return nil, false
	}

	return correctUser, true
}
