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
	normalizedTags := make([]string, 0, len(req.Tags))
	for _, t := range req.Tags {
		normalizedTags = append(normalizedTags, strings.ToLower(strings.TrimSpace(t)))
	}
	criteria := domain.SearchCriteria{
		Query: strings.TrimSpace(req.Query),
		Limit: 100,
	}
	places, err := uc.repo.Search(criteria)
	if err != nil {
		return SearchResponse{}, err
	}
	return SearchResponse{Places: places}, nil
}
