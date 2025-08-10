package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"src/internal/database"
	"src/internal/modules/places/application"
	"src/internal/modules/places/domain"
	osm "src/internal/modules/places/infrastructure/osm"
	pg "src/internal/modules/places/infrastructure/postgres"
	"src/internal/pkg/httpx"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	repo := pg.NewPlaceRepository(database.GormDB())

	r.Get("/", httpx.Endpoint(func(r *http.Request) (int, any, error) {
		// Backwards-compatible behavior: internal only, no external provider
		resp, err := application.NewSearchUseCase(repo).Execute(application.SearchRequest{})
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		dto := SearchResponseDTO{Places: toPlaceDTOs(resp.Places)}
		return http.StatusOK, dto, nil
	}))

	r.Get("/search", httpx.Endpoint(func(r *http.Request) (int, any, error) {
		// Extended search with provider/osm, geo, radius, pet type
		query := r.URL.Query().Get("q")
		tags := r.URL.Query()["tag"]
		provider := r.URL.Query().Get("provider")
		latPtr, lngPtr := parseLatLng(r)
		radiusPtr := parseRadius(r)
		petTypePtr := parsePetType(r)

		providers := map[string]domain.ExternalPlacesProvider{}
		providers["osm"] = osm.NewProvider("")
		uc := application.NewSearchAggregatedUseCase(repo, providers, nil)

		resp, err := uc.Execute(application.ExtendedSearchRequest{
			Query:    query,
			Tags:     tags,
			PetType:  petTypePtr,
			Lat:      latPtr,
			Lng:      lngPtr,
			RadiusM:  radiusPtr,
			Limit:    50,
			Provider: provider,
		})
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		dto := SearchResponseDTO{Places: toPlaceDTOs(resp.Places)}
		return http.StatusOK, dto, nil
	}))

	return r
}

// --- helpers ---
func parseLatLng(r *http.Request) (*float64, *float64) {
	q := r.URL.Query()
	latStr, lngStr := q.Get("lat"), q.Get("lng")
	if latStr == "" || lngStr == "" {
		return nil, nil
	}
	lat, err1 := strconv.ParseFloat(latStr, 64)
	lng, err2 := strconv.ParseFloat(lngStr, 64)
	if err1 != nil || err2 != nil {
		return nil, nil
	}
	return &lat, &lng
}

func parseRadius(r *http.Request) *int {
	q := r.URL.Query().Get("radius_m")
	if q == "" {
		return nil
	}
	if v, err := strconv.Atoi(q); err == nil {
		return &v
	}
	return nil
}

func parsePetType(r *http.Request) *domain.PetType {
	v := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("pet_type")))
	if v == "" {
		return nil
	}
	pt := domain.PetType(v)
	return &pt
}
