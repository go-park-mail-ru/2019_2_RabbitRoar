package usecase

import (
	"errors"
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

func (useCase* packUseCase) Delete(ID int, caller models.User) error {
	p, err := useCase.repo.GetByID(ID)
	if err != nil {
		return err
	}
	if p.Author != caller.ID {
		return errors.New("only creator can delete")
	}
	return useCase.repo.Delete(p.ID)
}

func (useCase* packUseCase) Update(pack *models.Pack, caller models.User) error {
	panic("implement me")
}

func (useCase packUseCase) GetByID(ID int, caller models.User) (*models.Pack, error) {
	//TODO: implement if user not played that pack or it marked as public offline or it is created by user
	p, err := useCase.repo.GetByID(ID)
	if err != nil {
		return nil, err
	}

	if p.Author != caller.ID {
		return nil, errors.New("not allowed")
	}

	return p, nil
}

func (useCase *packUseCase) FetchOffline(caller models.User) ([]int, error) {
	return useCase.repo.FetchOffline(caller)
}

func (useCase *packUseCase) FetchOfflinePublic() ([]int, error) {
	return useCase.repo.FetchOfflinePublic()
}

func (useCase* packUseCase) FetchOrderedByRating(desc bool, page, pageSize int) ([]models.Pack, error) {
	return useCase.repo.FetchOrderedByRating(desc, page, pageSize)
}

func (useCase* packUseCase) FetchByTags(tags string, page, pageSize int) ([]models.Pack, error) {
	return useCase.repo.FetchByTags(tags, page, pageSize)
}
