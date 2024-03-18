package transaction

import "github.com/google/uuid"

type AddMoneyInput struct {
	UserID uuid.UUID
	Amount float64
}

type AddMoneyOutput struct {
	Balance float64
}

type TransferMoneyInput struct {
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
	Amount     float64
}
