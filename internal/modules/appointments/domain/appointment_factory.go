package domain

import (
    "errors"
    "time"
)

var (
    ErrBookingInPast     = errors.New("cannot book in the past")
    ErrTooShortLeadTime  = errors.New("lead time too short")
    ErrInvalidPetID      = errors.New("invalid pet id")
    ErrInvalidSpecialist = errors.New("invalid specialist id")
)

// NewAppointment validates basic invariants and constructs an Appointment aggregate.
// leadTime is the minimal duration between now and start.
func NewAppointment(now time.Time, leadTime time.Duration, petID, specialistID string, slot TimeSlot) (Appointment, error) {
    if petID == "" {
        return Appointment{}, ErrInvalidPetID
    }
    if specialistID == "" {
        return Appointment{}, ErrInvalidSpecialist
    }
    if !now.Before(slot.Start) {
        return Appointment{}, ErrBookingInPast
    }
    if slot.Start.Sub(now) < leadTime {
        return Appointment{}, ErrTooShortLeadTime
    }
    return Appointment{
        PetID:        petID,
        SpecialistID: specialistID,
        StartTime:    slot.Start,
        EndTime:      slot.End,
        Status:       AppointmentStatusBooked,
    }, nil
}


