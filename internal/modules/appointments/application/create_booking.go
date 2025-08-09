package application

import (
	"context"
	"time"

	"src/internal/modules/appointments/domain"
	cport "src/internal/modules/shared/domain/clock"
	idport "src/internal/modules/shared/domain/idgen"
	txport "src/internal/modules/shared/domain/tx"
)

type CreateBookingRequest struct {
	PetID        string
	SpecialistID string
	StartTime    time.Time
	EndTime      time.Time
}

type CreateBookingResponse struct {
	Appointment domain.Appointment
}

type CreateBookingUseCase struct {
	repo         domain.AppointmentRepository
	idGen        idport.Port
	clock        cport.Port
	tx           txport.Manager
	availability domain.AvailabilityService
	policy       domain.BookingPolicy
	minLeadTime  time.Duration
}

func NewCreateBookingUseCase(repo domain.AppointmentRepository, idGen idport.Port, clock cport.Port, tx txport.Manager, availability domain.AvailabilityService, policy domain.BookingPolicy) *CreateBookingUseCase {
	return &CreateBookingUseCase{repo: repo, idGen: idGen, clock: clock, tx: tx, availability: availability, policy: policy, minLeadTime: 30 * time.Minute}
}

func (uc *CreateBookingUseCase) Execute(req CreateBookingRequest) (CreateBookingResponse, error) {
	var out CreateBookingResponse
	err := uc.tx.WithinTransaction(context.TODO(), func(_ context.Context) error {
		slot, err := domain.NewTimeSlot(req.StartTime, req.EndTime)
		if err != nil {
			return err
		}
		a, err := domain.NewAppointment(uc.clock.Now(), uc.minLeadTime, req.PetID, req.SpecialistID, slot)
		if err != nil {
			return err
		}
		if err := uc.availability.EnsureAvailable(req.SpecialistID, slot); err != nil {
			return err
		}
		if err := uc.policy.EnsureEligible(req.PetID, slot, uc.clock.Now()); err != nil {
			return err
		}
		a.ID = uc.idGen.NewID("apt")
		saved, err := uc.repo.Create(a)
		if err != nil {
			return err
		}
		out = CreateBookingResponse{Appointment: saved}
		return nil
	})
	if err != nil {
		return CreateBookingResponse{}, err
	}
	_ = uc.clock // reserved for future validation/time-based logic
	return out, nil
}
