package domain

// SearchCriteria captures filters for searching places across providers and internal DB.
type SearchCriteria struct {
	Query   string
	PetType *PetType
	Center  *GeoPoint
	RadiusM *int
	Limit   int
}

type PlaceRepository interface {
	Search(criteria SearchCriteria) ([]Place, error)
}

// ExternalPlacesProvider provides read-only access to third-party places (e.g., OSM/Google).
type ExternalPlacesProvider interface {
	ProviderName() string
	Search(criteria SearchCriteria) ([]Place, error)
}
