package domain

type PlaceRepository interface {
    Search(query string, tags []string) ([]Place, error)
}


