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
	normalizedTags := make([]string, 0, len(criteria.Tags))
	for _, t := range criteria.Tags {
		normalizedTags = append(normalizedTags, strings.ToLower(strings.TrimSpace(t)))
	}

	for _, p := range r.items {
		if criteria.Query != "" {
			if !strings.Contains(strings.ToLower(p.Name), strings.ToLower(criteria.Query)) &&
				!strings.Contains(strings.ToLower(p.Address), strings.ToLower(criteria.Query)) {
				continue
			}
		}
		if len(normalizedTags) > 0 {
			if !containsAll(stringsToLower(p.Tags), normalizedTags) {
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

func stringsToLower(xs []string) []string {
	ys := make([]string, 0, len(xs))
	for _, x := range xs {
		ys = append(ys, strings.ToLower(x))
	}
	return ys
}

func containsAll(haystack []string, needles []string) bool {
	for _, n := range needles {
		found := false
		for _, h := range haystack {
			if h == n {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
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
