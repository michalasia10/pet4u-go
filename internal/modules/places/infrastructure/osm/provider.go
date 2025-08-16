package osm

import (
	"context"
	"net/url"
	"time"

	"src/internal/cache"
	"src/internal/modules/places/domain"

	"github.com/redis/go-redis/v9"
)

type Provider struct {
	client   *OverpassClient
	cache    PlaceCache
	cacheTTL time.Duration
}

type ProviderOption func(*Provider)

func WithCache(cache PlaceCache, ttl time.Duration) ProviderOption {
	return func(p *Provider) {
		p.cache = cache
		p.cacheTTL = ttl
	}
}

func newProvider(overpassURL string, opts ...ProviderOption) *Provider {
	p := &Provider{
		client:   NewOverpassClient(overpassURL),
		cacheTTL: 12 * time.Hour,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func NewProviderWithInMemoryCache(overpassURL string, maxCacheSize int, ttl time.Duration) *Provider {
	cacheImpl := cache.NewInMemoryCache[[]domain.Place](maxCacheSize)
	return newProvider(overpassURL, WithCache(cacheImpl, ttl))
}

func NewProviderWithRedisCache(overpassURL string, redisClient *redis.Client, ttl time.Duration) *Provider {
	cacheImpl := cache.NewRedisCache[[]domain.Place](redisClient, "osm:places")
	return newProvider(overpassURL, WithCache(cacheImpl, ttl))
}

func NewProviderWithoutCache(overpassURL string) *Provider {
	return newProvider(overpassURL)
}

func (p *Provider) ProviderName() string { return "osm" }

func (p *Provider) Search(criteria domain.SearchCriteria) ([]domain.Place, error) {
	if criteria.Center == nil || criteria.RadiusM == nil {
		return []domain.Place{}, nil
	}
	cacheKey := GenerateCacheKey(criteria)

	if p.cache != nil {
		if places, found, err := p.cache.Get(context.Background(), cacheKey); err == nil && found {
			return places, nil
		}
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

	places := mapOverpassToPlaces(p.ProviderName(), payload)

	if p.cache != nil {
		p.cache.Set(context.Background(), cacheKey, places, p.cacheTTL)
	}

	return places, nil
}
