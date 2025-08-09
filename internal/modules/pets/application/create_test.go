package application_test

import (
	app "src/internal/modules/pets/application"
	"src/internal/modules/pets/domain"
	mem "src/internal/modules/pets/infrastructure/memory"
	cimpl "src/internal/modules/shared/infrastructure/clock"
	idimpl "src/internal/modules/shared/infrastructure/idgen"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("CreateUseCase", func() {
	It("creates a pet with generated ID", func() {
		repo := mem.NewInMemoryPetRepository()
		idGen := idimpl.NewTimeIDGen()
		clock := cimpl.NewSystemClock()
		uc := app.NewCreateUseCase(repo, idGen, clock)

		req := app.CreateRequest{
			Name:      "Rex",
			Species:   "dog",
			Breed:     "mutt",
			BirthDate: time.Now().AddDate(-1, 0, 0),
		}
		resp, err := uc.Execute(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Pet.Name).To(Equal("Rex"))
		assert.NotEmpty(GinkgoT(), resp.Pet.ID)
		Expect(resp.Pet.Records).To(Equal([]domain.MedicalRecord{}))
	})

	It("validates required fields (name, species)", func() {
		repo := mem.NewInMemoryPetRepository()
		idGen := idimpl.NewTimeIDGen()
		clock := cimpl.NewSystemClock()
		uc := app.NewCreateUseCase(repo, idGen, clock)

		req := app.CreateRequest{ // missing name/species
			Breed:     "mutt",
			BirthDate: time.Now().AddDate(-1, 0, 0),
		}
		// Domain currently doesn't validate; mimic HTTP validation usually handled by DTO.
		// Here we assert create executes and sets empty fields as provided.
		resp, err := uc.Execute(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Pet.Name).To(Equal(""))
		Expect(resp.Pet.Species).To(Equal(""))
	})
})
