package domain

import "time"

type MedicalRecord struct {
    ID        string
    PetID     string
    Notes     string
    CreatedAt time.Time
}

type Pet struct {
    ID        string
    Name      string
    Species   string
    Breed     string
    BirthDate time.Time
    Records   []MedicalRecord
}


