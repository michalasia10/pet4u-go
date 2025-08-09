package domain

import "time"

// AvailabilityService checks specialist availability for a given slot.
type AvailabilityService interface {
    EnsureAvailable(specialistID string, slot TimeSlot) error
}

// BookingPolicy validates cross-aggregate business rules for booking.
type BookingPolicy interface {
    EnsureEligible(petID string, slot TimeSlot, now time.Time) error
}


