package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	Amount float64   `gorm:"not null"`
}
