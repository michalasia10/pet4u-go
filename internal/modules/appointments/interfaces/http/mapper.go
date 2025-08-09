package http

import "src/internal/modules/appointments/domain"

func toAppointmentDTO(a domain.Appointment) AppointmentDTO {
    return AppointmentDTO{
        ID:           a.ID,
        PetID:        a.PetID,
        SpecialistID: a.SpecialistID,
        StartTime:    a.StartTime,
        EndTime:      a.EndTime,
        Status:       string(a.Status),
    }
}

func toAppointmentDTOs(items []domain.Appointment) []AppointmentDTO {
    out := make([]AppointmentDTO, 0, len(items))
    for _, it := range items {
        out = append(out, toAppointmentDTO(it))
    }
    return out
}


