package memory

import (
    "errors"
    "time"

    "src/internal/modules/appointments/domain"
)

var (
    ErrSlotTaken        = errors.New("slot overlaps existing appointment")
    ErrPetNotEligible   = errors.New("pet not eligible for booking")
)

// InMemoryAvailabilityService checks overlap against the in-memory repository.
type InMemoryAvailabilityService struct { Repo *InMemoryAppointmentRepository }

func NewInMemoryAvailabilityService(repo *InMemoryAppointmentRepository) *InMemoryAvailabilityService {
    return &InMemoryAvailabilityService{Repo: repo}
}

func (s *InMemoryAvailabilityService) EnsureAvailable(specialistID string, slot domain.TimeSlot) error {
    items, _ := s.Repo.List()
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

// InMemoryBookingPolicy is a stub policy (always eligible, or simple rule).
type InMemoryBookingPolicy struct{}

func NewInMemoryBookingPolicy() *InMemoryBookingPolicy { return &InMemoryBookingPolicy{} }

func (InMemoryBookingPolicy) EnsureEligible(petID string, slot domain.TimeSlot, now time.Time) error {
    // Example simple rule: allow all for now
    _ = petID; _ = slot; _ = now
    return nil
}


