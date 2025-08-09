package application

import "src/internal/modules/pets/domain"

type GetPetResponse struct {
    Pet domain.Pet `json:"pet"`
}

type GetPetUseCase struct {
    repo domain.PetRepository
}

func NewGetPetUseCase(repo domain.PetRepository) *GetPetUseCase {
    return &GetPetUseCase{repo: repo}
}

func (uc *GetPetUseCase) Execute(id string) (GetPetResponse, error) {
    pet, err := uc.repo.GetByID(id)
    if err != nil {
        return GetPetResponse{}, err
    }
    return GetPetResponse{Pet: pet}, nil
}


