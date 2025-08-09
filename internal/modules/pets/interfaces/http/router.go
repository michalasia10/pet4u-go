package http

import (
    "net/http"

    "github.com/go-chi/chi/v5"

    "src/internal/modules/pets/application"
    mem "src/internal/modules/pets/infrastructure/memory"
    "src/internal/pkg/httpx"
    cimpl "src/internal/modules/shared/infrastructure/clock"
    idimpl "src/internal/modules/shared/infrastructure/idgen"
)

func NewRouter() chi.Router {
    r := chi.NewRouter()
    repo := mem.NewInMemoryPetRepository()

    idGen := idimpl.NewTimeIDGen()
    clock := cimpl.NewSystemClock()
    createUC := application.NewCreatePetUseCase(repo, idGen, clock)
    getUC := application.NewGetPetUseCase(repo)

    r.Post("/", httpx.EndpointJSON[CreatePetRequestDTO](func(_ *http.Request, body CreatePetRequestDTO) (int, any, error) {
        if err := httpx.ValidateTags(body); err != nil {
            return http.StatusUnprocessableEntity, nil, err
        }
        resp, err := createUC.Execute(application.CreatePetRequest{
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


