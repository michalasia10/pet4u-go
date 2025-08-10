package domain

// GeoPoint represents a WGS84 latitude/longitude pair.
type GeoPoint struct {
	Lat float64
	Lng float64
}

// PetType is a coarse-grained animal kind used for filtering/search.
type PetType string

const (
	PetDog   PetType = "dog"
	PetCat   PetType = "cat"
	PetOther PetType = "other"
)

// SourceRef identifies an external provider reference for a place.
type SourceRef struct {
	Provider string
	ID       string
}

type Place struct {
	ID            string
	Name          string
	Address       string
	Location      GeoPoint
	Tags          []string
	PetTypes      []PetType
	IsPetFriendly bool
	Sources       []SourceRef
}
