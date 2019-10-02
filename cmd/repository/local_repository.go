package repository

import (
	"errors"

	"../entity"
	"github.com/google/uuid"
)

type LocalRepository struct {
	users       []entity.User
	sessions    map[uuid.UUID]entity.Session
	lastUserUID int64
}

func (repo LocalRepository) SessionGetUser(sessionId uuid.UUID) (entity.User, error) {
	if session, success := repo.sessions[sessionId]; success {
		return session.User, nil
	}

	return entity.User{}, errors.New("Session not found")
}

func (repo *LocalRepository) SessionCreate(user entity.User) (uuid.UUID, error) {
	newUUID, err := uuid.NewUUID()

	repo.sessions[newUUID] = entity.Session{
		Uuid: newUUID,
		User: user,
	}

	return newUUID, err
}

func (repo *LocalRepository) SessionDestroy(sessionId uuid.UUID) {
	delete(repo.sessions, sessionId)
}

func (repo LocalRepository) UserGetByName(name string) (entity.User, error) {
	for _, user := range repo.users {
		if user.Username == name {
			return user, nil
		}
	}
	return entity.User{}, ErrNotFound
}

func (repo *LocalRepository) UserCreate(name, password, email string) (entity.User, error) {
	_, err := repo.UserGetByName(name)
	if err != ErrNotFound {
		return entity.User{}, ErrConflict
	}

	repo.lastUserUID++
	user := entity.User{
		UID:      repo.lastUserUID,
		Username: name,
		Password: password,
		Email:    email,
		Rating:   0,
	}
	repo.users = append(repo.users, user)

	return user, nil
}

var Data LocalRepository

func init() {
	Data = LocalRepository{}
}
