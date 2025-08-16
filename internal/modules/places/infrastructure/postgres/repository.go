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

func (r *PlaceRepository) Search(criteria domain.SearchCriteria) ([]domain.Place, error) {
	tx := r.db.WithContext(context.Background()).Model(&Place{})
	if criteria.Query != "" {
		like := "%" + criteria.Query + "%"
		tx = tx.Where("LOWER(name) LIKE LOWER(?) OR LOWER(address) LIKE LOWER(?)", like, like)
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
