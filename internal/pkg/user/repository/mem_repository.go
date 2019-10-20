package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/user"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
)

type memUserRepository struct {
	users       []models.User
	lastUserUID int64
}

var ErrUserNotFound = errors.New("error user not found")
var ErrUserConflict = errors.New("error conflict")

func NewMemUserRepository() user.Repository {
	return &memUserRepository{}
}

func (repo *memUserRepository) GetByName(name string) (*models.User, error) {
	for _, u := range repo.users {
		if u.Username == name {
			return &u, nil
		}
	}
	return nil, ErrUserNotFound
}

func (repo *memUserRepository) GetByID(userID int) (*models.User, error) {
	if len(repo.users) < userID {
		return nil, ErrUserNotFound
	}

	u := repo.users[userID]
	return &u, nil
}

func (repo *memUserRepository) Update(updatedUser models.User) error {
	var err error = nil

	defer func() {
		if r := recover(); r != nil {
			err = ErrUserNotFound
		}
	}()

	repo.users[updatedUser.UID] = updatedUser

	return err
}

func (repo *memUserRepository) Create(u models.User) (*models.User, error) {
	_, err := repo.GetByName(u.Username)
	if err != ErrUserNotFound {
		return nil, ErrUserConflict
	}

	u.UID = repo.lastUserUID
	repo.lastUserUID++
	repo.users = append(repo.users, u)

	return &u, nil
}
