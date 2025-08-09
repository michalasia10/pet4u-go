package application

import (
	"time"

	"src/internal/modules/pets/domain"
	cport "src/internal/modules/shared/domain/clock"
	idport "src/internal/modules/shared/domain/idgen"
)

type CreateRequest struct {
	Name      string
	Species   string
	Breed     string
	BirthDate time.Time
}

type CreateResponse struct {
	Pet domain.Pet
}

type CreateUseCase struct {
	repo  domain.PetRepository
	idGen idport.Port
	clock cport.Port
}

func NewCreateUseCase(repo domain.PetRepository, idGen idport.Port, clock cport.Port) *CreateUseCase {
	return &CreateUseCase{repo: repo, idGen: idGen, clock: clock}
}

func (uc *CreateUseCase) Execute(req CreateRequest) (CreateResponse, error) {
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
		return CreateResponse{}, err
	}
	_ = uc.clock // reserved for future time-based logic
	return CreateResponse{Pet: saved}, nil
}
