package osm

import (
	"fmt"
	"src/internal/modules/places/domain"
)

func mapOverpassToPlaces(providerName string, payload overpassResponse) []domain.Place {
	places := make([]domain.Place, 0, len(payload.Elements))
	for _, el := range payload.Elements {
		name := el.Tags["name"]
		if name == "" {
			continue
		}
		address := composeAddress(el.Tags)
		ptypes := []domain.PetType{}
		if el.Tags["amenity"] == "veterinary" || el.Tags["shop"] == "pet" || el.Tags["leisure"] == "dog_park" {
			ptypes = append(ptypes, domain.PetDog, domain.PetCat)
		}
		place := domain.Place{
			ID:            "",
			Name:          name,
			Address:       address,
			Location:      domain.GeoPoint{Lat: el.LatOrCenterLat(), Lng: el.LonOrCenterLon()},
			Tags:          deriveTags(el.Tags),
			PetTypes:      ptypes,
			IsPetFriendly: true,
			Sources:       []domain.SourceRef{{Provider: providerName, ID: el.ElementID()}},
		}
		places = append(places, place)
	}
	return places
}

func composeAddress(tags map[string]string) string {
	parts := []string{}
	if v := tags["addr:street"]; v != "" {
		parts = append(parts, v)
	}
	if v := tags["addr:housenumber"]; v != "" {
		if len(parts) > 0 {
			parts[len(parts)-1] = fmt.Sprintf("%s %s", parts[len(parts)-1], v)
		} else {
			parts = append(parts, v)
		}
	}
	if v := tags["addr:city"]; v != "" {
		parts = append(parts, v)
	}
	return stringsJoin(parts, ", ")
}

func deriveTags(tags map[string]string) []string {
	out := []string{}
	for k, v := range tags {
		switch {
		case k == "amenity" && (v == "veterinary" || v == "drinking_water"):
			out = append(out, v)
		case k == "shop" && v == "pet":
			out = append(out, "pet_shop")
		case k == "leisure" && v == "dog_park":
			out = append(out, v)
		}
	}
	return out
}

// lightweight string helpers to avoid extra imports
func stringsJoin(elems []string, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0]
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}
	b := make([]byte, n)
	bp := copy(b, elems[0])
	for _, s := range elems[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}
	return string(b)
}
