package domain

import "time"

type AppointmentStatus string

const (
    AppointmentStatusPending  AppointmentStatus = "pending"
    AppointmentStatusBooked   AppointmentStatus = "booked"
    AppointmentStatusCancelled AppointmentStatus = "cancelled"
)

type Appointment struct {
    ID           string
    PetID        string
    SpecialistID string
    StartTime    time.Time
    EndTime      time.Time
    Status       AppointmentStatus
}


