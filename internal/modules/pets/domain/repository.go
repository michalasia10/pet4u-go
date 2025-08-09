package domain

type PetRepository interface {
    Create(p Pet) (Pet, error)
    GetByID(id string) (Pet, error)
    Update(p Pet) (Pet, error)
    List() ([]Pet, error)
}


