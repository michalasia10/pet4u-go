package postgres

import (
	"gorm.io/datatypes"
)

type Place struct {
	ID            string         `gorm:"primaryKey;type:uuid"`
	Name          string         `gorm:"type:varchar(255);not null;index"`
	Address       string         `gorm:"type:varchar(255);not null"`
	Tags          datatypes.JSON `gorm:"type:jsonb;default:'[]'::jsonb;not null"`
	IsPetFriendly bool           `gorm:"not null;default:true"`
}

func (Place) TableName() string { return "places" }
