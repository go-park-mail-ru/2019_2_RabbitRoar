package models

//go:generate easyjson -all

import "github.com/google/uuid"

type Game struct {
	UUID            uuid.UUID
	Name            string `json:"name"`
	PlayersCapacity int    `json:"playersCapacity"`
	PlayersJoined   int    `json:"playersJoined"`
	PackID          int    `json:"pack"`
	PackName        string `json:"packName"`
}
