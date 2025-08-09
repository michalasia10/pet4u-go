package application_test

import (
	app "src/internal/modules/pets/application"
	"src/internal/modules/pets/domain"
	mem "src/internal/modules/pets/infrastructure/memory"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetUseCase", func() {
	It("returns a pet by id", func() {
		repo := mem.NewInMemoryPetRepository()
		repo.Create(domain.Pet{ID: "p1", Name: "Rex"})
		uc := app.NewGetUseCase(repo)
		resp, err := uc.Execute("p1")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Pet.Name).To(Equal("Rex"))
	})

	It("propagates error when pet not found", func() {
		repo := mem.NewInMemoryPetRepository()
		uc := app.NewGetUseCase(repo)
		_, err := uc.Execute("missing")
		Expect(err).To(HaveOccurred())
	})
})
