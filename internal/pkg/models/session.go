package models

//go:generate easyjson -all

import (
	"github.com/google/uuid"
)

type Session struct {
	UUID uuid.UUID
	User User
}
