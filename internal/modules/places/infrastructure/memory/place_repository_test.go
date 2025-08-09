package memory_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"src/internal/modules/places/domain"
	. "src/internal/modules/places/infrastructure/memory"
)

var _ = Describe("InMemoryPlaceRepository", func() {
	It("searches by name/address and tags", func() {
		seed := []domain.Place{
			{ID: "1", Name: "Cafe Paws", Address: "123 Bark St", Tags: []string{"cafe", "wifi"}},
			{ID: "2", Name: "Happy Park", Address: "Green Ave", Tags: []string{"park"}},
		}
		repo := NewInMemoryPlaceRepository(seed)

		resp, err := repo.Search("park", nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp).To(HaveLen(1))

		resp, err = repo.Search("", []string{"WIFI"})
		Expect(err).ToNot(HaveOccurred())
		Expect(resp).To(HaveLen(1))
	})
})
