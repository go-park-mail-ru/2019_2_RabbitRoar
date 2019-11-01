package pack

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"

type Repository interface {
	GetByID(ID int) (*models.Pack, error)
	GetQuestions(pack models.Pack) (*[]models.Question, error)
	FetchByRating(asc bool, page, pageSize int) (*[]models.Pack, error)
	FetchByTags(tags string, page, pageSize int) (*[]models.Pack, error)
	Create(pack *models.Pack) error
	Update(pack *models.Pack) error
	Delete(ID int) error
}
