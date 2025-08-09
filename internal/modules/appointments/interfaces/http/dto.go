package http

import "time"

type CreateRequestDTO struct {
	PetID        string    `json:"pet_id" validate:"required"`
	SpecialistID string    `json:"specialist_id" validate:"required"`
	StartTime    time.Time `json:"start_time" validate:"required"`
	EndTime      time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
}

type AppointmentDTO struct {
	ID           string    `json:"id"`
	PetID        string    `json:"pet_id"`
	SpecialistID string    `json:"specialist_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       string    `json:"status"`
}

type ListResponseDTO struct {
	Appointments []AppointmentDTO `json:"appointments"`
}
