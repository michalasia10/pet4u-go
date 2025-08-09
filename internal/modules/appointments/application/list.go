package application

import "src/internal/modules/appointments/domain"

type ListResponse struct {
	Appointments []domain.Appointment
}

type ListUseCase struct {
	repo domain.AppointmentRepository
}

func NewListUseCase(repo domain.AppointmentRepository) *ListUseCase {
	return &ListUseCase{repo: repo}
}

func (uc *ListUseCase) Execute() (ListResponse, error) {
	items, err := uc.repo.List()
	if err != nil {
		return ListResponse{}, err
	}
	return ListResponse{Appointments: items}, nil
}
