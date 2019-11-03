package models

import (
	"github.com/google/uuid"
)

type Session struct {
	UUID uuid.UUID
	User User
}
