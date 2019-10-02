package entity

import (
	"github.com/google/uuid"
)

type Session struct {
	Uuid uuid.UUID
	User User
}
