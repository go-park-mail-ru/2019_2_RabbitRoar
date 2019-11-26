package pack

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"

type UseCase interface {
	Create(pack *models.Pack, caller models.User) error
	Update(pack *models.Pack) error
	Delete(ID int) error
	Played(packID, userID int) bool
	GetByID(ID int) (*models.Pack, error)
	FetchOffline(caller models.User) ([]int, error)
	FetchOfflinePublic() ([]int, error)
	FetchOfflineAuthor(caller models.User) ([]int, error)
	FetchByAuthor(author models.User, desc bool, page, pageSize int) ([]models.Pack, error)
	FetchOrderedByRating(desc bool, page, pageSize int) ([]models.Pack, error)
	FetchByTags(tags string, page, pageSize int) ([]models.Pack, error)
}
