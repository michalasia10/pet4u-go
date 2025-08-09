package http

import "time"

type CreateRequestDTO struct {
	Name      string    `json:"name" validate:"required"`
	Species   string    `json:"species"`
	Breed     string    `json:"breed"`
	BirthDate time.Time `json:"birth_date"`
}

// FieldErrors is a simple map-based validation error container.
type FieldErrors map[string]string

func (e FieldErrors) Error() string { return "validation failed" }

func (d CreateRequestDTO) Validate() error {
	errs := FieldErrors{}
	if d.Name == "" {
		errs["name"] = "required"
	}
	if d.Species == "" {
		errs["species"] = "required"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type RecordDTO struct {
	ID        string    `json:"id"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type PetDTO struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Species   string      `json:"species"`
	Breed     string      `json:"breed"`
	BirthDate time.Time   `json:"birth_date"`
	Records   []RecordDTO `json:"records"`
}
