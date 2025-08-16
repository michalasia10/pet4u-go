package osm_test

import (
	"context"
	"testing"
	"time"

	"src/internal/cache"
	"src/internal/modules/places/domain"
	"src/internal/modules/places/infrastructure/osm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	var (
		placeCache osm.PlaceCache
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Describe("InMemoryCache", func() {
		BeforeEach(func() {
			placeCache = cache.NewInMemoryCache[[]domain.Place](10) // small max size for testing
		})

		It("should return miss for non-existent key", func() {
			places, found, err := placeCache.Get(ctx, "non-existent")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
			Expect(places).To(BeNil())
		})

		It("should store and retrieve data", func() {
			places := []domain.Place{
				{ID: "1", Name: "Test Place 1"},
				{ID: "2", Name: "Test Place 2"},
			}

			err := placeCache.Set(ctx, "test-key", places, time.Hour)
			Expect(err).NotTo(HaveOccurred())

			retrieved, found, err := placeCache.Get(ctx, "test-key")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(retrieved).To(HaveLen(2))
			Expect(retrieved[0].Name).To(Equal("Test Place 1"))
			Expect(retrieved[1].Name).To(Equal("Test Place 2"))
		})

		It("should handle TTL expiration", func() {
			places := []domain.Place{{ID: "1", Name: "Test Place"}}

			err := placeCache.Set(ctx, "test-key", places, 1*time.Millisecond)
			Expect(err).NotTo(HaveOccurred())

			// Wait for expiration
			time.Sleep(10 * time.Millisecond)

			retrieved, found, err := placeCache.Get(ctx, "test-key")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
			Expect(retrieved).To(BeNil())
		})

		It("should clear all data", func() {
			places := []domain.Place{{ID: "1", Name: "Test Place"}}

			err := placeCache.Set(ctx, "test-key-1", places, time.Hour)
			Expect(err).NotTo(HaveOccurred())
			err = placeCache.Set(ctx, "test-key-2", places, time.Hour)
			Expect(err).NotTo(HaveOccurred())

			err = placeCache.Clear(ctx)
			Expect(err).NotTo(HaveOccurred())

			_, found1, _ := placeCache.Get(ctx, "test-key-1")
			_, found2, _ := placeCache.Get(ctx, "test-key-2")
			Expect(found1).To(BeFalse())
			Expect(found2).To(BeFalse())
		})

		It("should track cache stats", func() {
			places := []domain.Place{{ID: "1", Name: "Test Place"}}

			// Set some data
			placeCache.Set(ctx, "test-key", places, time.Hour)

			// Hit
			placeCache.Get(ctx, "test-key")

			// Miss
			placeCache.Get(ctx, "non-existent")

			stats := placeCache.Stats()
			Expect(stats.Hits).To(Equal(int64(1)))
			Expect(stats.Misses).To(Equal(int64(1)))
			Expect(stats.Size).To(Equal(int64(1)))
		})

		It("should evict oldest entry when max size reached", func() {
			smallCache := cache.NewInMemoryCache[[]domain.Place](2) // max 2 entries
			places := []domain.Place{{ID: "1", Name: "Test Place"}}

			// Fill cache to capacity
			smallCache.Set(ctx, "key1", places, time.Hour)
			smallCache.Set(ctx, "key2", places, time.Hour)

			// Both should be present
			_, found1, _ := smallCache.Get(ctx, "key1")
			_, found2, _ := smallCache.Get(ctx, "key2")
			Expect(found1).To(BeTrue())
			Expect(found2).To(BeTrue())

			// Add third entry - should evict oldest
			smallCache.Set(ctx, "key3", places, time.Hour)

			// key1 should be evicted, key2 and key3 should remain
			_, found1, _ = smallCache.Get(ctx, "key1")
			_, found2, _ = smallCache.Get(ctx, "key2")
			_, found3, _ := smallCache.Get(ctx, "key3")
			Expect(found1).To(BeFalse())
			Expect(found2).To(BeTrue())
			Expect(found3).To(BeTrue())
		})
	})

	Describe("GenerateCacheKey", func() {
		It("should generate consistent keys for same criteria", func() {
			criteria := domain.SearchCriteria{
				Query:   "veterinary",
				Center:  &domain.GeoPoint{Lat: 52.2297, Lng: 21.0122},
				RadiusM: &[]int{1000}[0],
			}

			key1 := osm.GenerateCacheKey(criteria)
			key2 := osm.GenerateCacheKey(criteria)

			Expect(key1).To(Equal(key2))
			Expect(key1).To(HavePrefix("osm:search:"))
		})

		It("should generate different keys for different criteria", func() {
			criteria1 := domain.SearchCriteria{
				Query:   "veterinary",
				Center:  &domain.GeoPoint{Lat: 52.2297, Lng: 21.0122},
				RadiusM: &[]int{1000}[0],
			}

			criteria2 := domain.SearchCriteria{
				Query:   "park",
				Center:  &domain.GeoPoint{Lat: 52.2297, Lng: 21.0122},
				RadiusM: &[]int{1000}[0],
			}

			key1 := osm.GenerateCacheKey(criteria1)
			key2 := osm.GenerateCacheKey(criteria2)

			Expect(key1).NotTo(Equal(key2))
		})

		It("should generate different keys for different locations", func() {
			criteria1 := domain.SearchCriteria{
				Query:   "veterinary",
				Center:  &domain.GeoPoint{Lat: 52.2297, Lng: 21.0122},
				RadiusM: &[]int{1000}[0],
			}

			criteria2 := domain.SearchCriteria{
				Query:   "veterinary",
				Center:  &domain.GeoPoint{Lat: 51.7592, Lng: 19.4560},
				RadiusM: &[]int{1000}[0],
			}

			key1 := osm.GenerateCacheKey(criteria1)
			key2 := osm.GenerateCacheKey(criteria2)

			Expect(key1).NotTo(Equal(key2))
		})
	})
})

func TestCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OSM Cache Suite")
}
