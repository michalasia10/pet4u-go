package http

type SearchRequestDTO struct {
	Query string   `json:"query"`
	Tags  []string `json:"tags"`
}

type PlaceDTO struct {
	ID            string                     `json:"id"`
	Name          string                     `json:"name"`
	Address       string                     `json:"address"`
	Location      struct{ Lat, Lng float64 } `json:"location"`
	Tags          []string                   `json:"tags"`
	PetTypes      []string                   `json:"pet_types"`
	IsPetFriendly bool                       `json:"is_pet_friendly"`
	Source        string                     `json:"source"`
}

type SearchResponseDTO struct {
	Places []PlaceDTO `json:"places"`
}
