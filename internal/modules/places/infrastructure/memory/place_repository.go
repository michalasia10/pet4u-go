package memory

import (
	"strings"

	"src/internal/modules/places/domain"
)

type InMemoryPlaceRepository struct {
	items []domain.Place
}

func NewInMemoryPlaceRepository(seed []domain.Place) *InMemoryPlaceRepository {
	return &InMemoryPlaceRepository{items: seed}
}

func (r *InMemoryPlaceRepository) Search(criteria domain.SearchCriteria) ([]domain.Place, error) {
	var result []domain.Place

	for _, p := range r.items {
		if criteria.Query != "" {
			if !strings.Contains(strings.ToLower(p.Name), strings.ToLower(criteria.Query)) &&
				!strings.Contains(strings.ToLower(p.Address), strings.ToLower(criteria.Query)) {
				continue
			}
		}
		// PetType filter
		if criteria.PetType != nil {
			if !placeSupportsPetType(p, *criteria.PetType) {
				continue
			}
		}
		result = append(result, p)
	}
	// Limit if requested
	if criteria.Limit > 0 && len(result) > criteria.Limit {
		result = result[:criteria.Limit]
	}
	return result, nil
}

func placeSupportsPetType(p domain.Place, petType domain.PetType) bool {
	if len(p.PetTypes) == 0 {
		// If unspecified, assume pet-friendly places are okay for any pet
		return p.IsPetFriendly
	}
	for _, pt := range p.PetTypes {
		if pt == petType {
			return true
		}
	}
	return false
}
