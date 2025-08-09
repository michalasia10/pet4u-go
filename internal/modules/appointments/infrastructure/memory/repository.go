package memory

import "src/internal/modules/appointments/domain"

type InMemoryAppointmentRepository struct {
    items []domain.Appointment
}

func NewInMemoryAppointmentRepository() *InMemoryAppointmentRepository {
    return &InMemoryAppointmentRepository{items: make([]domain.Appointment, 0)}
}

func (r *InMemoryAppointmentRepository) Create(a domain.Appointment) (domain.Appointment, error) {
    r.items = append(r.items, a)
    return a, nil
}

func (r *InMemoryAppointmentRepository) List() ([]domain.Appointment, error) {
    return append([]domain.Appointment(nil), r.items...), nil
}


