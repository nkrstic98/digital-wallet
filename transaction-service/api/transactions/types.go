package transactions

import "github.com/google/uuid"

type AddMoneyRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Amount float64   `json:"amount"`
}

type AddMoneyResponse struct {
	Balance float64 `json:"balance"`
}

type TransferMoneyRequest struct {
	FromUserID uuid.UUID `json:"from_user_id"`
	ToUserID   uuid.UUID `json:"to_user_id"`
	Amount     float64   `json:"amount"`
}
