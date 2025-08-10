package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"src/internal/database"
	"src/internal/modules/places/application"
	pg "src/internal/modules/places/infrastructure/postgres"
	"src/internal/pkg/httpx"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	repo := pg.NewPlaceRepository(database.GormDB())
	uc := application.NewSearchUseCase(repo)

	r.Get("/", httpx.Endpoint(func(r *http.Request) (int, any, error) {
		resp, err := uc.Execute(application.SearchRequest{})
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		dto := SearchResponseDTO{Places: toPlaceDTOs(resp.Places)}
		return http.StatusOK, dto, nil
	}))

	r.Get("/search", httpx.Endpoint(func(r *http.Request) (int, any, error) {
		query := r.URL.Query().Get("q")
		tags := r.URL.Query()["tag"]
		resp, err := uc.Execute(application.SearchRequest{Query: query, Tags: tags})
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		dto := SearchResponseDTO{Places: toPlaceDTOs(resp.Places)}
		return http.StatusOK, dto, nil
	}))

	return r
}
