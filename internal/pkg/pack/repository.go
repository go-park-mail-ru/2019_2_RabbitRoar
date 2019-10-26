package pack

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"

type Repository interface {
	GetByID(int) (*models.Pack, error)
	GetQuestions(pack models.Pack) (*[]models.Question, error)
	Create(pack models.Pack) (*models.Pack, error)
	Update(pack models.Pack) (*models.Pack, error)
	Delete(int) error
}
