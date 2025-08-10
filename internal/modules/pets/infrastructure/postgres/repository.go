package postgres

import (
	"context"

	"gorm.io/gorm"

	"src/internal/modules/pets/domain"
)

type PetRepository struct {
	db *gorm.DB
}

func NewPetRepository(db *gorm.DB) *PetRepository {
	return &PetRepository{db: db}
}

func toRecordPet(p domain.Pet) Pet {
	return Pet{
		ID:        p.ID,
		Name:      p.Name,
		Species:   p.Species,
		Breed:     p.Breed,
		BirthDate: p.BirthDate,
	}
}

func toDomainPet(r Pet) domain.Pet {
	return domain.Pet{
		ID:        r.ID,
		Name:      r.Name,
		Species:   r.Species,
		Breed:     r.Breed,
		BirthDate: r.BirthDate,
		Records:   nil,
	}
}

func (r *PetRepository) Create(p domain.Pet) (domain.Pet, error) {
	rec := toRecordPet(p)
	if err := r.db.WithContext(context.Background()).Create(&rec).Error; err != nil {
		return domain.Pet{}, err
	}
	return toDomainPet(rec), nil
}

func (r *PetRepository) GetByID(id string) (domain.Pet, error) {
	var rec Pet
	if err := r.db.WithContext(context.Background()).First(&rec, "id = ?", id).Error; err != nil {
		return domain.Pet{}, err
	}
	return toDomainPet(rec), nil
}

func (r *PetRepository) Update(p domain.Pet) (domain.Pet, error) {
	rec := toRecordPet(p)
	if err := r.db.WithContext(context.Background()).Model(&Pet{}).Where("id = ?", p.ID).Updates(rec).Error; err != nil {
		return domain.Pet{}, err
	}
	return toDomainPet(rec), nil
}

func (r *PetRepository) List() ([]domain.Pet, error) {
	var rows []Pet
	if err := r.db.WithContext(context.Background()).Order("name asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Pet, 0, len(rows))
	for _, row := range rows {
		out = append(out, toDomainPet(row))
	}
	return out, nil
}
