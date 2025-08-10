package osm

import (
	"context"
	"net/url"

	"src/internal/modules/places/domain"
)

// Provider implements domain.ExternalPlacesProvider using OpenStreetMap Overpass API.
// It orchestrates query building, calling the client and mapping to domain.
type Provider struct {
	client *OverpassClient
}

func NewProvider(overpassURL string) *Provider {
	return &Provider{client: NewOverpassClient(overpassURL)}
}

func (p *Provider) ProviderName() string { return "osm" }

func (p *Provider) Search(criteria domain.SearchCriteria) ([]domain.Place, error) {
	if criteria.Center == nil || criteria.RadiusM == nil {
		return []domain.Place{}, nil
	}

	filters := []string{"amenity=veterinary", "shop=pet", "amenity=park", "leisure=dog_park"}
	if criteria.PetType != nil && *criteria.PetType == domain.PetDog {
		filters = append(filters, "amenity=drinking_water")
	}

	nameRegex := ""
	if criteria.Query != "" {
		// best effort sanitization; Overpass supports regex
		_ = url.QueryEscape(criteria.Query)
		nameRegex = criteria.Query
	}

	ql := NewQueryBuilder().
		Around(*criteria.RadiusM, criteria.Center.Lat, criteria.Center.Lng).
		WithFilters(filters...).
		NameRegex(nameRegex).
		Build()

	payload, err := p.client.Execute(context.Background(), ql)
	if err != nil {
		return nil, err
	}
	return mapOverpassToPlaces(p.ProviderName(), payload), nil
}
