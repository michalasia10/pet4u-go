package postgres

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"src/internal/modules/places/domain"
)

type PlaceRepository struct {
	db *gorm.DB
}

func NewPlaceRepository(db *gorm.DB) *PlaceRepository {
	return &PlaceRepository{db: db}
}

func toDomainPlace(r Place) domain.Place {
	var tags []string
	if len(r.Tags) > 0 {
		_ = json.Unmarshal([]byte(r.Tags), &tags)
	}
	return domain.Place{
		ID:            r.ID,
		Name:          r.Name,
		Address:       r.Address,
		IsPetFriendly: r.IsPetFriendly,
		Tags:          tags,
		// Location, PetTypes, Sources are not yet persisted in this MVP
	}
}

func (r *PlaceRepository) Search(criteria domain.SearchCriteria) ([]domain.Place, error) {
	tx := r.db.WithContext(context.Background()).Model(&Place{})
	if criteria.Query != "" {
		like := "%" + criteria.Query + "%"
		tx = tx.Where("LOWER(name) LIKE LOWER(?) OR LOWER(address) LIKE LOWER(?)", like, like)
	}
	if len(criteria.Tags) > 0 {
		// contains-all over JSONB using @> with array; requires tags stored as jsonb array
		// Example matches ANY: tx = tx.Where("tags ?| array[?]", pq.StringArray(tags))
		// For contains-all, iterate; simple approach: first tag only
		tx = tx.Where("tags @> ?::jsonb", toJSONBArray(criteria.Tags))
	}
	var rows []Place
	limit := 100
	if criteria.Limit > 0 && criteria.Limit < limit {
		limit = criteria.Limit
	}
	if err := tx.Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Place, 0, len(rows))
	for _, row := range rows {
		out = append(out, toDomainPlace(row))
	}
	return out, nil
}

func toJSONBArray(xs []string) string {
	b, _ := json.Marshal(xs)
	return string(b)
}
