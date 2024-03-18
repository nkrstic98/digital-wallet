package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Email     string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
