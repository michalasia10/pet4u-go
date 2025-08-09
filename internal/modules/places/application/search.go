package application

import (
	"strings"

	"src/internal/modules/places/domain"
)

type SearchRequest struct {
	Query string
	Tags  []string
}

type SearchResponse struct {
	Places []domain.Place
}

type SearchUseCase struct {
	repo domain.PlaceRepository
}

func NewSearchUseCase(repo domain.PlaceRepository) *SearchUseCase {
	return &SearchUseCase{repo: repo}
}

func (uc *SearchUseCase) Execute(req SearchRequest) (SearchResponse, error) {
	places, err := uc.repo.Search(strings.TrimSpace(req.Query), req.Tags)
	if err != nil {
		return SearchResponse{}, err
	}
	return SearchResponse{Places: places}, nil
}
