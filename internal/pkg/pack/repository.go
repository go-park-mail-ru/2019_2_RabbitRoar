package pack

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/pkg/errors"
)

var(
	ErrRepoNotFound = errors.New("pack with such ID not found")
	ErrRepoCorrupted = errors.New("pack corrupted")
)

type Repository interface {
	Create(pack *models.Pack) error
	Update(pack *models.Pack) error
	Delete(ID int) error
	GetByID(ID int) (*models.Pack, error)
	Played(packID, userID int) (bool, error)
	FetchOffline(caller models.User) ([]int, error)
	FetchOfflineAuthor(caller models.User) ([]int, error)
	FetchOfflinePublic() ([]int, error)
	FetchByAuthor(u models.User, desc bool, page, pageSize int) ([]models.Pack, error)
	FetchOrderedByRating(desc bool, page, pageSize int) ([]models.Pack, error)
	FetchByTags(tags string, desc bool, page, pageSize int) ([]models.Pack, error)
}
