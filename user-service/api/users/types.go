package users

import (
	"github.com/google/uuid"
	"time"
)

type CreateRequest struct {
	Email string `json:"email"`
}

type CreateResponse struct {
	UserId    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type GetUserBalanceResponse struct {
	Email   string  `json:"email"`
	Balance float64 `json:"balance"`
}
