package usecase

import (
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/pack"
)

type packUseCase struct {
	repo pack.Repository
}

func NewUserUseCase(repo pack.Repository) pack.UseCase {
	return &packUseCase{
		repo: repo,
	}
}

func (useCase* packUseCase) Create(p *models.Pack, caller models.User) error {
	p.Author = caller.ID
	return useCase.repo.Create(p)
}

func (useCase* packUseCase) Delete(ID int) error {
	return useCase.repo.Delete(ID)
}

func (useCase *packUseCase) Update(pack *models.Pack, caller models.User) error {
	panic("implement me")
}

func (useCase *packUseCase) Played(packID, userID int) bool {
	played, err := useCase.repo.Played(packID, userID)
	if err != nil {
		//TODO: log error here
		return false
	}
	return played
}

func (useCase *packUseCase) GetByID(ID int) (*models.Pack, error) {
	return useCase.repo.GetByID(ID)
}

func (useCase *packUseCase) FetchOffline(caller models.User) ([]int, error) {
	return useCase.repo.FetchOffline(caller)
}

func (useCase *packUseCase) FetchOfflineAuthor(caller models.User) ([]int, error) {
	return useCase.repo.FetchOfflineAuthor(caller)
}

func (useCase *packUseCase) FetchOfflinePublic() ([]int, error) {
	return useCase.repo.FetchOfflinePublic()
}

func (useCase* packUseCase) FetchOrderedByRating(desc bool, page, pageSize int) ([]models.Pack, error) {
	return useCase.repo.FetchOrderedByRating(desc, page, pageSize)
}

func (useCase* packUseCase) FetchByAuthor(author models.User, desc bool, page, pageSize int) ([]models.Pack, error) {
	return useCase.repo.FetchByAuthor(author, true, page, pageSize)
}

func (useCase* packUseCase) FetchByTags(tags string, page, pageSize int) ([]models.Pack, error) {
	return useCase.repo.FetchByTags(tags, true, page, pageSize)
}
