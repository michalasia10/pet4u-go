package postgres

import (
	"time"
)

type Appointment struct {
	ID           string    `gorm:"primaryKey;type:uuid"`
	PetID        string    `gorm:"type:uuid;not null;index"`
	SpecialistID string    `gorm:"type:uuid;not null;index"`
	StartTime    time.Time `gorm:"not null;index"`
	EndTime      time.Time `gorm:"not null"`
	Status       string    `gorm:"type:varchar(32);not null;index"`
}

func (Appointment) TableName() string { return "appointments" }
