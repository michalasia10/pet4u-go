package http

import "src/internal/modules/places/domain"

func toPlaceDTO(p domain.Place) PlaceDTO {
	return PlaceDTO{
		ID:            p.ID,
		Name:          p.Name,
		Address:       p.Address,
		Location:      struct{ Lat, Lng float64 }{Lat: p.Location.Lat, Lng: p.Location.Lng},
		Tags:          append([]string(nil), p.Tags...),
		PetTypes:      petTypesToStrings(p.PetTypes),
		IsPetFriendly: p.IsPetFriendly,
		Source:        firstSource(p.Sources),
	}
}

func toPlaceDTOs(xs []domain.Place) []PlaceDTO {
	out := make([]PlaceDTO, 0, len(xs))
	for _, p := range xs {
		out = append(out, toPlaceDTO(p))
	}
	return out
}

func petTypesToStrings(xs []domain.PetType) []string {
	out := make([]string, 0, len(xs))
	for _, x := range xs {
		out = append(out, string(x))
	}
	return out
}

func firstSource(xs []domain.SourceRef) string {
	if len(xs) == 0 {
		return "internal"
	}
	return xs[0].Provider
}
