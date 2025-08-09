package application

import (
    "strings"

    "src/internal/modules/places/domain"
)

type SearchPlacesRequest struct {
    Query string
    Tags  []string
}

type SearchPlacesResponse struct {
    Places []domain.Place `json:"places"`
}

type SearchPlacesUseCase struct {
    repo domain.PlaceRepository
}

func NewSearchPlacesUseCase(repo domain.PlaceRepository) *SearchPlacesUseCase {
    return &SearchPlacesUseCase{repo: repo}
}

func (uc *SearchPlacesUseCase) Execute(req SearchPlacesRequest) (SearchPlacesResponse, error) {
    places, err := uc.repo.Search(strings.TrimSpace(req.Query), req.Tags)
    if err != nil {
        return SearchPlacesResponse{}, err
    }
    return SearchPlacesResponse{Places: places}, nil
}


