package usecase

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session"
	"github.com/google/uuid"
)

type sessionUseCase struct {
	repository session.Repository
}

func NewSessionUseCase(repository session.Repository) session.UseCase {
	return &sessionUseCase{
		repository: repository,
	}
}

func (uc sessionUseCase) GetUserByUUID(u uuid.UUID) (*models.User, error) {
	return uc.repository.GetUser(u)
}

func (uc sessionUseCase) GetUserByStringUUID(su string) (*models.User, error) {
	UUID, err := uuid.Parse(su)
	if err != nil {
		return nil, err
	}
	return uc.repository.GetUser(UUID)
}

func (uc sessionUseCase) Create(u models.User) (*uuid.UUID, error) {
	return uc.repository.Create(u)
}

func (uc sessionUseCase) Destroy(sessionId uuid.UUID) {
	uc.repository.Destroy(sessionId)
}
