package http

import "src/internal/modules/places/domain"

func toPlaceDTO(p domain.Place) PlaceDTO {
    return PlaceDTO{
        ID:            p.ID,
        Name:          p.Name,
        Address:       p.Address,
        Tags:          append([]string(nil), p.Tags...),
        IsPetFriendly: p.IsPetFriendly,
    }
}

func toPlaceDTOs(xs []domain.Place) []PlaceDTO {
    out := make([]PlaceDTO, 0, len(xs))
    for _, p := range xs {
        out = append(out, toPlaceDTO(p))
    }
    return out
}


