package postgres

import (
	"time"
)

// Pet is the persistence model for pets.
type Pet struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Species   string    `gorm:"type:varchar(100);not null"`
	Breed     string    `gorm:"type:varchar(100);not null"`
	BirthDate time.Time `gorm:"not null"`
}

func (Pet) TableName() string { return "pets" }

// MedicalRecord is the persistence model for pet medical records.
type MedicalRecord struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	PetID     string    `gorm:"type:uuid;not null;index"`
	Notes     string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null;index"`
}

func (MedicalRecord) TableName() string { return "medical_records" }
