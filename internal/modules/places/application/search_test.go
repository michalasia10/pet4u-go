package application_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	app "src/internal/modules/places/application"
	"src/internal/modules/places/domain"
	mem "src/internal/modules/places/infrastructure/memory"
	"testing"
)

func TestPlacesApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Places Application Suite")
}

var _ = Describe("SearchUseCase", func() {
	var (
		repo *mem.InMemoryPlaceRepository
		uc   *app.SearchUseCase
	)

	BeforeEach(func() {
		seed := []domain.Place{
			{ID: "1", Name: "Cafe Paws", Address: "123 Bark St", Tags: []string{"cafe", "wifi"}, IsPetFriendly: true},
			{ID: "2", Name: "Happy Park", Address: "Green Ave", Tags: []string{"park"}, IsPetFriendly: true},
			{ID: "3", Name: "Cat Corner", Address: "Feline Rd", Tags: []string{"cafe"}, IsPetFriendly: true},
		}
		repo = mem.NewInMemoryPlaceRepository(seed)
		uc = app.NewSearchUseCase(repo)
	})

	It("returns all when no filters provided", func() {
		resp, err := uc.Execute(app.SearchRequest{})
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Places).To(HaveLen(3))
	})

	It("filters by query (case-insensitive, name/address)", func() {
		resp, err := uc.Execute(app.SearchRequest{Query: "park"})
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Places).To(HaveLen(1))
		Expect(resp.Places[0].Name).To(Equal("Happy Park"))
	})

	It("filters by tags (must contain all)", func() {
		resp, err := uc.Execute(app.SearchRequest{Tags: []string{"cafe"}})
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Places).To(HaveLen(2))
	})

	It("normalizes tags and query (trim/lower)", func() {
		resp, err := uc.Execute(app.SearchRequest{Query: "  CaFe  ", Tags: []string{"  WiFi "}})
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Places).ToNot(BeEmpty())
	})

	It("plays nice with testify require for additional checks", func() {
		t := GinkgoT()
		resp, err := uc.Execute(app.SearchRequest{Query: "cafe"})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(resp.Places), 1)
	})
})
