package application

import (
    "time"

    "src/internal/modules/pets/domain"
    cport "src/internal/modules/shared/domain/clock"
    idport "src/internal/modules/shared/domain/idgen"
)

type CreatePetRequest struct {
    Name      string
    Species   string
    Breed     string
    BirthDate time.Time
}

type CreatePetResponse struct {
    Pet domain.Pet `json:"pet"`
}

type CreatePetUseCase struct {
    repo domain.PetRepository
    idGen idport.Port
    clock cport.Port
}

func NewCreatePetUseCase(repo domain.PetRepository, idGen idport.Port, clock cport.Port) *CreatePetUseCase {
    return &CreatePetUseCase{repo: repo, idGen: idGen, clock: clock}
}

func (uc *CreatePetUseCase) Execute(req CreatePetRequest) (CreatePetResponse, error) {
    p := domain.Pet{
        ID:        uc.idGen.NewID("pet"),
        Name:      req.Name,
        Species:   req.Species,
        Breed:     req.Breed,
        BirthDate: req.BirthDate,
        Records:   []domain.MedicalRecord{},
    }
    saved, err := uc.repo.Create(p)
    if err != nil {
        return CreatePetResponse{}, err
    }
    _ = uc.clock // reserved for future time-based logic
    return CreatePetResponse{Pet: saved}, nil
}


