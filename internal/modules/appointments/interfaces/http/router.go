package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"src/internal/modules/appointments/application"
	mem "src/internal/modules/appointments/infrastructure/memory"
	cimpl "src/internal/modules/shared/infrastructure/clock"
	idimpl "src/internal/modules/shared/infrastructure/idgen"
	tximpl "src/internal/modules/shared/infrastructure/tx"
	"src/internal/pkg/httpx"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	repo := mem.NewInMemoryAppointmentRepository()
	idGen := idimpl.NewTimeIDGen()
	clock := cimpl.NewSystemClock()
	tx := tximpl.NewNoopManager()
	availability := mem.NewInMemoryAvailabilityService(repo)
	policy := mem.NewInMemoryBookingPolicy()
	createUC := application.NewCreateBookingUseCase(repo, idGen, clock, tx, availability, policy)
	listUC := application.NewListBookingsUseCase(repo)

	r.Get("/", httpx.Endpoint(func(r *http.Request) (int, any, error) {
		resp, err := listUC.Execute()
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		dto := ListBookingsResponseDTO{Appointments: toAppointmentDTOs(resp.Appointments)}
		return http.StatusOK, dto, nil
	}))

	r.Post("/", httpx.EndpointJSON[CreateBookingRequestDTO](func(_ *http.Request, body CreateBookingRequestDTO) (int, any, error) {
		if err := httpx.ValidateTags(body); err != nil {
			return http.StatusUnprocessableEntity, nil, err
		}
		resp, err := createUC.Execute(application.CreateBookingRequest{
			PetID:        body.PetID,
			SpecialistID: body.SpecialistID,
			StartTime:    body.StartTime,
			EndTime:      body.EndTime,
		})
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		dto := toAppointmentDTO(resp.Appointment)
		return http.StatusCreated, dto, nil
	}))

	return r
}
