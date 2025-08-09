package http

import "src/internal/modules/pets/domain"

func toPetDTO(p domain.Pet) PetDTO {
    return PetDTO{
        ID:        p.ID,
        Name:      p.Name,
        Species:   p.Species,
        Breed:     p.Breed,
        BirthDate: p.BirthDate,
        Records:   toMedicalRecordDTOs(p.Records),
    }
}

func toMedicalRecordDTOs(xs []domain.MedicalRecord) []MedicalRecordDTO {
    out := make([]MedicalRecordDTO, 0, len(xs))
    for _, r := range xs {
        out = append(out, MedicalRecordDTO{ID: r.ID, Notes: r.Notes, CreatedAt: r.CreatedAt})
    }
    return out
}


