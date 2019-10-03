package repository

import (
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserCreateGet(t *testing.T) {
	var rep LocalRepository
	userCreated, errCreated := rep.UserCreate("iamfirstuser", "fuckpassword", "test@mail.ru")
	assert.Nil(t, errCreated)
	user, err := rep.UserGetByName("iamfirstuser")
	assert.Nil(t, err)
	assert.Equal(t, userCreated, user)
}

func TestUserCreateUIDIncrement(t *testing.T) {
	var rep LocalRepository
	user1, _ := rep.UserCreate("1", "1234", "test@mail.ru")
	user2, _ := rep.UserCreate("2", "1234", "test@mail.ru")
	assert.NotEqual(t, user1.UID, user2.UID)
}

func TestCreateUserConflict(t *testing.T) {
	var rep LocalRepository
	rep.UserCreate("1", "1234", "test@mail.ru")
	_, err := rep.UserCreate("1", "1234", "test@mail.ru")
	assert.Equal(t, ErrConflict, err)
}

func TestSessionCreate(t *testing.T) {
	rep := LocalRepository{}
	rep.sessions = make(map[uuid.UUID]entity.Session)
	user1, _ := rep.UserCreate("1", "1234", "test@mail.ru")

	userUUID, err := rep.SessionCreate(user1)
	assert.Nil(t, err)

	userBySession, err := rep.SessionGetUser(userUUID)
	assert.Equal(t, user1, userBySession)
}

func TestSessionDestroy(t *testing.T) {
	rep := LocalRepository{}
	rep.sessions = make(map[uuid.UUID]entity.Session)
	user1, _ := rep.UserCreate("1", "1234", "test@mail.ru")

	userUUID, err := rep.SessionCreate(user1)
	assert.Nil(t, err)

	rep.SessionDestroy(userUUID)

	_, err = rep.SessionGetUser(userUUID)
	assert.Equal(t, errors.New("Session not found"), err)
}
