package memory

import (
    "errors"

    "src/internal/modules/pets/domain"
)

type InMemoryPetRepository struct {
    items map[string]domain.Pet
}

func NewInMemoryPetRepository() *InMemoryPetRepository {
    return &InMemoryPetRepository{items: make(map[string]domain.Pet)}
}

func (r *InMemoryPetRepository) Create(p domain.Pet) (domain.Pet, error) {
    r.items[p.ID] = p
    return p, nil
}

func (r *InMemoryPetRepository) GetByID(id string) (domain.Pet, error) {
    p, ok := r.items[id]
    if !ok {
        return domain.Pet{}, errors.New("not found")
    }
    return p, nil
}

func (r *InMemoryPetRepository) Update(p domain.Pet) (domain.Pet, error) {
    r.items[p.ID] = p
    return p, nil
}

func (r *InMemoryPetRepository) List() ([]domain.Pet, error) {
    result := make([]domain.Pet, 0, len(r.items))
    for _, p := range r.items {
        result = append(result, p)
    }
    return result, nil
}


