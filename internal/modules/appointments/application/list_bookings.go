package application

import "src/internal/modules/appointments/domain"

type ListBookingsResponse struct {
	Appointments []domain.Appointment
}

type ListBookingsUseCase struct {
	repo domain.AppointmentRepository
}

func NewListBookingsUseCase(repo domain.AppointmentRepository) *ListBookingsUseCase {
	return &ListBookingsUseCase{repo: repo}
}

func (uc *ListBookingsUseCase) Execute() (ListBookingsResponse, error) {
	items, err := uc.repo.List()
	if err != nil {
		return ListBookingsResponse{}, err
	}
	return ListBookingsResponse{Appointments: items}, nil
}
