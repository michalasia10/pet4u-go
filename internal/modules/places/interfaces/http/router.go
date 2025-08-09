package http

import (
    "net/http"

    "github.com/go-chi/chi/v5"

    "src/internal/modules/places/application"
    "src/internal/modules/places/domain"
    mem "src/internal/modules/places/infrastructure/memory"
    "src/internal/pkg/httpx"
)

func NewRouter() chi.Router {
    r := chi.NewRouter()

    seed := []domain.Place{
        {ID: "pl_1", Name: "Cafe Paws", Address: "123 Bark St", Tags: []string{"cafe", "wifi"}, IsPetFriendly: true},
        {ID: "pl_2", Name: "Happy Park", Address: "Green Ave", Tags: []string{"park"}, IsPetFriendly: true},
    }
    repo := mem.NewInMemoryPlaceRepository(seed)
    uc := application.NewSearchPlacesUseCase(repo)

    r.Get("/", httpx.Endpoint(func(r *http.Request) (int, any, error) {
        resp, err := uc.Execute(application.SearchPlacesRequest{})
        if err != nil {
            return http.StatusInternalServerError, nil, err
        }
        dto := SearchPlacesResponseDTO{Places: toPlaceDTOs(resp.Places)}
        return http.StatusOK, dto, nil
    }))

    r.Get("/search", httpx.Endpoint(func(r *http.Request) (int, any, error) {
        query := r.URL.Query().Get("q")
        tags := r.URL.Query()["tag"]
        resp, err := uc.Execute(application.SearchPlacesRequest{Query: query, Tags: tags})
        if err != nil {
            return http.StatusInternalServerError, nil, err
        }
        dto := SearchPlacesResponseDTO{Places: toPlaceDTOs(resp.Places)}
        return http.StatusOK, dto, nil
    }))

    return r
}

