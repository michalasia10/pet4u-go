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
	}
}

func (r *PlaceRepository) Search(query string, tags []string) ([]domain.Place, error) {
	tx := r.db.WithContext(context.Background()).Model(&Place{})
	if query != "" {
		like := "%" + query + "%"
		tx = tx.Where("LOWER(name) LIKE LOWER(?) OR LOWER(address) LIKE LOWER(?)", like, like)
	}
	if len(tags) > 0 {
		// contains-all over JSONB using @> with array; requires tags stored as jsonb array
		// Example matches ANY: tx = tx.Where("tags ?| array[?]", pq.StringArray(tags))
		// For contains-all, iterate; simple approach: first tag only
		tx = tx.Where("tags @> ?::jsonb", toJSONBArray(tags))
	}
	var rows []Place
	if err := tx.Limit(100).Find(&rows).Error; err != nil {
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
