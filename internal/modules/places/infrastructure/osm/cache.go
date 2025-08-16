package osm

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"src/internal/cache"
	"src/internal/modules/places/domain"
)

type PlaceCache = cache.Cache[[]domain.Place]

func GenerateCacheKey(criteria domain.SearchCriteria) string {
	data := struct {
		Query   string
		PetType *domain.PetType
		Center  *domain.GeoPoint
		RadiusM *int
	}{
		Query:   criteria.Query,
		PetType: criteria.PetType,
		Center:  criteria.Center,
		RadiusM: criteria.RadiusM,
	}

	jsonBytes, _ := json.Marshal(data)
	hash := md5.Sum(jsonBytes)
	return fmt.Sprintf("osm:search:%x", hash)
}
