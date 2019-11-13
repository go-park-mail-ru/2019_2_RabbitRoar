package models

import "github.com/google/uuid"

type Game struct {
	UUID            uuid.UUID
	Name            string
	PlayersCapacity int
	PlayersJoined   int
	Creator         int
	PackID          int
}
