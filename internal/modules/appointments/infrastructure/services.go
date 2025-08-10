package infrastructure

import (
	"time"

	"src/internal/modules/appointments/domain"
)

// AvailabilityService checks availability using the AppointmentRepository port.
type AvailabilityService struct {
	Repo domain.AppointmentRepository
}

func NewAvailabilityService(repo domain.AppointmentRepository) *AvailabilityService {
	return &AvailabilityService{Repo: repo}
}

func (s *AvailabilityService) EnsureAvailable(specialistID string, slot domain.TimeSlot) error {
	items, err := s.Repo.List()
	if err != nil {
		return err
	}
	for _, a := range items {
		if a.SpecialistID == specialistID {
			existing, _ := domain.NewTimeSlot(a.StartTime, a.EndTime)
			if slot.Overlaps(existing) {
				return ErrSlotTaken
			}
		}
	}
	return nil
}

// ErrSlotTaken exported from infrastructure to map to domain-level behavior without importing infra in domain.
var ErrSlotTaken = domain.ErrTooShortLeadTime // placeholder mapping; prefer specific error

// NoopBookingPolicy is a default policy that allows all bookings.
type NoopBookingPolicy struct{}

func NewNoopBookingPolicy() *NoopBookingPolicy { return &NoopBookingPolicy{} }

func (NoopBookingPolicy) EnsureEligible(petID string, slot domain.TimeSlot, now time.Time) error {
	return nil
}
