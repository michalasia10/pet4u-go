package http

type SearchPlacesRequestDTO struct {
    Query string   `json:"query"`
    Tags  []string `json:"tags"`
}

type PlaceDTO struct {
    ID            string   `json:"id"`
    Name          string   `json:"name"`
    Address       string   `json:"address"`
    Tags          []string `json:"tags"`
    IsPetFriendly bool     `json:"is_pet_friendly"`
}

type SearchPlacesResponseDTO struct {
    Places []PlaceDTO `json:"places"`
}


