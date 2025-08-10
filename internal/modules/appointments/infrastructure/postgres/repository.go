package postgres

import (
	"context"

	"gorm.io/gorm"

	"src/internal/modules/appointments/domain"
)

type AppointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func toRecordAppointment(a domain.Appointment) Appointment {
	return Appointment{
		ID:           a.ID,
		PetID:        a.PetID,
		SpecialistID: a.SpecialistID,
		StartTime:    a.StartTime,
		EndTime:      a.EndTime,
		Status:       string(a.Status),
	}
}

func toDomainAppointment(r Appointment) domain.Appointment {
	return domain.Appointment{
		ID:           r.ID,
		PetID:        r.PetID,
		SpecialistID: r.SpecialistID,
		StartTime:    r.StartTime,
		EndTime:      r.EndTime,
		Status:       domain.AppointmentStatus(r.Status),
	}
}

func (r *AppointmentRepository) Create(a domain.Appointment) (domain.Appointment, error) {
	rec := toRecordAppointment(a)
	if err := r.db.WithContext(context.Background()).Create(&rec).Error; err != nil {
		return domain.Appointment{}, err
	}
	return toDomainAppointment(rec), nil
}

func (r *AppointmentRepository) List() ([]domain.Appointment, error) {
	var rows []Appointment
	if err := r.db.WithContext(context.Background()).Order("start_time asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Appointment, 0, len(rows))
	for _, row := range rows {
		out = append(out, toDomainAppointment(row))
	}
	return out, nil
}
