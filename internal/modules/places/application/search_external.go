package application

import (
	"math"
	"sort"
	"strings"

	"src/internal/modules/places/domain"
)

// Extended request supporting external providers and geo/pet filters.
type ExtendedSearchRequest struct {
	Query    string
	PetType  *domain.PetType
	Lat      *float64
	Lng      *float64
	RadiusM  *int
	Limit    int
	Provider string // optional, e.g., "osm"; empty means default provider
}

type ExtendedSearchResponse struct {
	Places []domain.Place
}

type Taxonomy interface {
	KeywordsFor(pet domain.PetType) []string
}

// SimpleTaxonomy is a pragmatic in-memory taxonomy implementation.
type SimpleTaxonomy struct{}

func (SimpleTaxonomy) KeywordsFor(pet domain.PetType) []string {
	switch pet {
	case domain.PetDog:
		return []string{"dog", "dog park", "veterinary", "groomer", "pet store", "animal shelter"}
	case domain.PetCat:
		return []string{"cat", "veterinary", "groomer", "pet store", "animal shelter"}
	default:
		return []string{"veterinary", "pet store"}
	}
}

// SearchAggregatedUseCase orchestrates external provider and internal repository.
type SearchAggregatedUseCase struct {
	providers map[string]domain.ExternalPlacesProvider
	repo      domain.PlaceRepository
	taxonomy  Taxonomy
}

func NewSearchAggregatedUseCase(repo domain.PlaceRepository, providers map[string]domain.ExternalPlacesProvider, taxonomy Taxonomy) *SearchAggregatedUseCase {
	if taxonomy == nil {
		taxonomy = SimpleTaxonomy{}
	}
	if providers == nil {
		providers = map[string]domain.ExternalPlacesProvider{}
	}
	return &SearchAggregatedUseCase{repo: repo, providers: providers, taxonomy: taxonomy}
}

func (uc *SearchAggregatedUseCase) Execute(req ExtendedSearchRequest) (ExtendedSearchResponse, error) {
	criteria := domain.SearchCriteria{
		Query: strings.TrimSpace(req.Query),
		Limit: clampPositive(req.Limit, 1, 50),
	}
	if req.PetType != nil {
		criteria.PetType = req.PetType
		// Extend query with taxonomy keywords to increase recall on providers
		keywords := uc.taxonomy.KeywordsFor(*req.PetType)
		if len(keywords) > 0 && criteria.Query == "" {
			criteria.Query = strings.Join(keywords, ", ")
		}
	}
	if req.Lat != nil && req.Lng != nil {
		criteria.Center = &domain.GeoPoint{Lat: *req.Lat, Lng: *req.Lng}
	}
	if req.RadiusM != nil {
		r := clampPositive(*req.RadiusM, 100, 20000)
		criteria.RadiusM = &r
	}

	// External first; if it fails, degrade to internal
	var externalPlaces []domain.Place
	providerName := strings.TrimSpace(strings.ToLower(req.Provider))
	if providerName == "" {
		providerName = "osm"
	}
	if provider, ok := uc.providers[providerName]; ok && provider != nil {
		if xs, err := provider.Search(criteria); err == nil {
			externalPlaces = xs
		}
	}

	internalPlaces, err := uc.repo.Search(criteria)
	if err != nil {
		return ExtendedSearchResponse{}, err
	}

	merged := mergeAndDedupe(externalPlaces, internalPlaces)
	// Sort by distance if geo provided; otherwise keep provider-first ordering
	if criteria.Center != nil {
		center := *criteria.Center
		sort.SliceStable(merged, func(i, j int) bool {
			di := distanceMeters(center, merged[i].Location)
			dj := distanceMeters(center, merged[j].Location)
			return di < dj
		})
	}

	if criteria.Limit > 0 && len(merged) > criteria.Limit {
		merged = merged[:criteria.Limit]
	}
	return ExtendedSearchResponse{Places: merged}, nil
}

func clampPositive(v, min, max int) int {
	if v <= 0 {
		return max
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func mergeAndDedupe(primary, secondary []domain.Place) []domain.Place {
	// Mark seen by provider-id and by name+geo proximity
	seenExternal := map[string]struct{}{}
	out := make([]domain.Place, 0, len(primary)+len(secondary))
	for _, p := range primary {
		for _, s := range p.Sources {
			key := s.Provider + ":" + s.ID
			seenExternal[key] = struct{}{}
		}
		out = append(out, p)
	}
	for _, p := range secondary {
		if isDuplicateBySource(p, seenExternal) {
			continue
		}
		if containsNearDuplicate(out, p) {
			continue
		}
		out = append(out, p)
	}
	return out
}

func isDuplicateBySource(p domain.Place, seen map[string]struct{}) bool {
	for _, s := range p.Sources {
		key := s.Provider + ":" + s.ID
		if _, ok := seen[key]; ok {
			return true
		}
	}
	return false
}

func containsNearDuplicate(existing []domain.Place, candidate domain.Place) bool {
	nameNorm := normalizeName(candidate.Name)
	for _, e := range existing {
		if normalizeName(e.Name) != nameNorm {
			continue
		}
		// treat as near-duplicate when distance < 100m
		if distanceMeters(e.Location, candidate.Location) < 100 {
			return true
		}
	}
	return false
}

func normalizeName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	replacers := []string{" ltd", " sp. z o.o.", ",", ".", "'", "\""}
	for _, r := range replacers {
		s = strings.ReplaceAll(s, r, "")
	}
	s = strings.Join(strings.Fields(s), " ")
	return s
}

// Haversine distance in meters
func distanceMeters(a, b domain.GeoPoint) float64 {
	const earthRadiusM = 6371000.0
	toRad := func(d float64) float64 { return d * math.Pi / 180 }
	dLat := toRad(b.Lat - a.Lat)
	dLng := toRad(b.Lng - a.Lng)
	lat1 := toRad(a.Lat)
	lat2 := toRad(b.Lat)
	sinDLat := math.Sin(dLat / 2)
	sinDLng := math.Sin(dLng / 2)
	h := sinDLat*sinDLat + math.Cos(lat1)*math.Cos(lat2)*sinDLng*sinDLng
	c := 2 * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))
	return earthRadiusM * c
}
