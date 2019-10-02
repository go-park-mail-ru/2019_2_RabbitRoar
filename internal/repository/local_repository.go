package repository

import (
	"../entity"
)

type LocalRepository struct {
	users       []entity.User
	lastUserUID int64
}

func (repo LocalRepository) UserGetByName(name string) (entity.User, error) {
	for _, user := range repo.users {
		if user.Username == name {
			return user, nil
		}
	}
	return entity.User{}, ErrNotFound
}

func (repo *LocalRepository) UserCreate(name, password string) (entity.User, error) {
	_, err := repo.UserGetByName(name)
	if err != ErrNotFound {
		return entity.User{}, ErrConflict
	}

	repo.lastUserUID++
	user := entity.User{
		UID:      repo.lastUserUID,
		Username: name,
		Password: password,
		Rating:   0,
	}
	repo.users = append(repo.users, user)

	return user, nil
}
