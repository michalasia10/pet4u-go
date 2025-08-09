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

func (r *InMemoryPlaceRepository) Search(query string, tags []string) ([]domain.Place, error) {
    var result []domain.Place
    normalizedTags := make([]string, 0, len(tags))
    for _, t := range tags {
        normalizedTags = append(normalizedTags, strings.ToLower(strings.TrimSpace(t)))
    }

    for _, p := range r.items {
        if query != "" {
            if !strings.Contains(strings.ToLower(p.Name), strings.ToLower(query)) &&
                !strings.Contains(strings.ToLower(p.Address), strings.ToLower(query)) {
                continue
            }
        }
        if len(normalizedTags) > 0 {
            if !containsAll(stringsToLower(p.Tags), normalizedTags) {
                continue
            }
        }
        result = append(result, p)
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


