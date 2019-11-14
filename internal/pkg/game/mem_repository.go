package game

import "github.com/google/uuid"

type MemRepository interface {
	JoinConnection(gameID uuid.UUID, conn PlayerConnection) error
}
