package models

import "github.com/google/uuid"

type GameState int

const (
	Pending = iota
	Started
)

type Game struct {
	UUID            uuid.UUID
	Name            string
	PlayersCapacity int
	PlayersJoined   int
	State           GameState
	Creator         int
}
