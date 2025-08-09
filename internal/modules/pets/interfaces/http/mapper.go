package http

import "src/internal/modules/pets/domain"

func toPetDTO(p domain.Pet) PetDTO {
	return PetDTO{
		ID:        p.ID,
		Name:      p.Name,
		Species:   p.Species,
		Breed:     p.Breed,
		BirthDate: p.BirthDate,
		Records:   toRecordDTOs(p.Records),
	}
}

func toRecordDTOs(xs []domain.MedicalRecord) []RecordDTO {
	out := make([]RecordDTO, 0, len(xs))
	for _, r := range xs {
		out = append(out, RecordDTO{ID: r.ID, Notes: r.Notes, CreatedAt: r.CreatedAt})
	}
	return out
}
