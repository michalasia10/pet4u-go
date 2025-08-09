package application

import "src/internal/modules/pets/domain"

type GetResponse struct {
	Pet domain.Pet
}

type GetUseCase struct {
	repo domain.PetRepository
}

func NewGetUseCase(repo domain.PetRepository) *GetUseCase {
	return &GetUseCase{repo: repo}
}

func (uc *GetUseCase) Execute(id string) (GetResponse, error) {
	pet, err := uc.repo.GetByID(id)
	if err != nil {
		return GetResponse{}, err
	}
	return GetResponse{Pet: pet}, nil
}
