package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"src/internal/database"
	"src/internal/modules/appointments/application"
	infra "src/internal/modules/appointments/infrastructure"
	pg "src/internal/modules/appointments/infrastructure/postgres"
	cimpl "src/internal/modules/shared/infrastructure/clock"
	idimpl "src/internal/modules/shared/infrastructure/idgen"
	tximpl "src/internal/modules/shared/infrastructure/tx"
	"src/internal/pkg/httpx"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	repo := pg.NewAppointmentRepository(database.GormDB())
	idGen := idimpl.NewUUIDGen()
	clock := cimpl.NewSystemClock()
	tx := tximpl.NewNoopManager()
	availability := infra.NewAvailabilityService(repo)
	policy := infra.NewNoopBookingPolicy()
	createUC := application.NewCreateUseCase(repo, idGen, clock, tx, availability, policy)
	listUC := application.NewListUseCase(repo)

	r.Get("/", httpx.Endpoint(func(r *http.Request) (int, any, error) {
		resp, err := listUC.Execute()
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		dto := ListResponseDTO{Appointments: toAppointmentDTOs(resp.Appointments)}
		return http.StatusOK, dto, nil
	}))

	r.Post("/", httpx.EndpointJSON[CreateRequestDTO](func(_ *http.Request, body CreateRequestDTO) (int, any, error) {
		if err := httpx.ValidateTags(body); err != nil {
			return http.StatusUnprocessableEntity, nil, err
		}
		resp, err := createUC.Execute(application.CreateRequest{
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
