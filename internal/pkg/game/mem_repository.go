package game

import "github.com/google/uuid"

type MemRepository interface {
	Create(gameID uuid.UUID, hostID int) error
	JoinConnection(gameID uuid.UUID, conn Connection) error
}
