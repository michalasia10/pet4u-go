package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"src/internal/database"
	"src/internal/modules/pets/application"
	pg "src/internal/modules/pets/infrastructure/postgres"
	cimpl "src/internal/modules/shared/infrastructure/clock"
	idimpl "src/internal/modules/shared/infrastructure/idgen"
	"src/internal/pkg/httpx"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	repo := pg.NewPetRepository(database.GormDB())

	idGen := idimpl.NewUUIDGen()
	clock := cimpl.NewSystemClock()
	createUC := application.NewCreateUseCase(repo, idGen, clock)
	getUC := application.NewGetUseCase(repo)

	r.Post("/", httpx.EndpointJSON[CreateRequestDTO](func(_ *http.Request, body CreateRequestDTO) (int, any, error) {
		if err := httpx.ValidateTags(body); err != nil {
			return http.StatusUnprocessableEntity, nil, err
		}
		resp, err := createUC.Execute(application.CreateRequest{
			Name:      body.Name,
			Species:   body.Species,
			Breed:     body.Breed,
			BirthDate: body.BirthDate,
		})
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		return http.StatusCreated, toPetDTO(resp.Pet), nil
	}))

	r.Get("/{id}", httpx.Endpoint(func(r *http.Request) (int, any, error) {
		id := chi.URLParam(r, "id")
		resp, err := getUC.Execute(id)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		return http.StatusOK, toPetDTO(resp.Pet), nil
	}))

	return r
}
