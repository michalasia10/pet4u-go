package domain

type AppointmentRepository interface {
    Create(a Appointment) (Appointment, error)
    List() ([]Appointment, error)
}


