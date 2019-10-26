package question

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"

type Repository interface {
	GetByID(int) (*models.Question, error)
	Create(question models.Question) error
	Update(question models.Question) error
	Delete(int) error
}
