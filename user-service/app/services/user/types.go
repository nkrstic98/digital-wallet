package user

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UserId    uuid.UUID
	Email     string
	CreatedAt time.Time
}

type UserCreatedEventRequest struct {
	UserId    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
