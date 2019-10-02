package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserCreateGet(t *testing.T) {
	var rep LocalRepository
	userCreated, errCreated := rep.UserCreate("iamfirstuser", "fuckpassword")
	assert.Nil(t, errCreated)
	user, err := rep.UserGetByName("iamfirstuser")
	assert.Nil(t, err)
	assert.Equal(t, userCreated, user)
}

func TestUserCreateUIDIncrement(t *testing.T) {
	var rep LocalRepository
	user1, _ := rep.UserCreate("1", "1234")
	user2, _ := rep.UserCreate("2", "1234")
	assert.NotEqual(t, user1.UID, user2.UID)
}

func TestCreateUserConflict(t *testing.T) {
	var rep LocalRepository
	rep.UserCreate("1", "1234")
	_, err := rep.UserCreate("1", "1234")
	assert.Equal(t, ErrConflict, err)
}
