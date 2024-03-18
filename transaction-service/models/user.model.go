package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Balance   float64   `gorm:"not null"`
	CreatedAt time.Time
}
