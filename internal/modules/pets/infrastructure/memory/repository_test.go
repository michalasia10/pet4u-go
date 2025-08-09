package memory_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"src/internal/modules/pets/domain"
	. "src/internal/modules/pets/infrastructure/memory"
)

var _ = Describe("InMemoryPetRepository", func() {
	var repo *InMemoryPetRepository

	BeforeEach(func() {
		repo = NewInMemoryPetRepository()
	})

	It("creates and retrieves pets", func() {
		p := domain.Pet{ID: "p1", Name: "Rex"}
		_, err := repo.Create(p)
		Expect(err).ToNot(HaveOccurred())

		got, err := repo.GetByID("p1")
		Expect(err).ToNot(HaveOccurred())
		Expect(got.Name).To(Equal("Rex"))
	})

	It("returns error on missing pet", func() {
		_, err := repo.GetByID("missing")
		Expect(err).To(HaveOccurred())
		Expect(errors.Is(err, err)).To(BeTrue())
	})

	It("updates and lists pets", func() {
		p := domain.Pet{ID: "p1", Name: "Rex"}
		repo.Create(p)

		p.Name = "Max"
		repo.Update(p)

		list, err := repo.List()
		Expect(err).ToNot(HaveOccurred())
		Expect(list).To(HaveLen(1))
		Expect(list[0].Name).To(Equal("Max"))
	})
})
