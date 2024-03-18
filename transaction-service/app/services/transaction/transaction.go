package transaction

import (
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	"transaction-service/models"
)

type ServiceImpl struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *ServiceImpl {
	return &ServiceImpl{
		db: db,
	}
}

func (service *ServiceImpl) AddMoney(input AddMoneyInput) (float64, error) {
	var balance float64

	err := service.db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		result := tx.Clauses(clause.Locking{
			Strength: clause.LockingStrengthUpdate,
		}).Where("id = ?", input.UserID).First(&user)
		if result.Error != nil {
			slog.Error(fmt.Sprintf("Error fetching user with id %v: %s", input.UserID, result.Error.Error()))
			return result.Error
		}

		user.Balance += input.Amount

		result = tx.Save(&user)
		if result.Error != nil {
			slog.Error(fmt.Sprintf("Error updating user with id %v: %s", input.UserID, result.Error.Error()))
			return result.Error
		}

		transaction := models.Transaction{
			UserID: input.UserID,
			Amount: input.Amount,
		}
		result = tx.Create(&transaction)
		if result.Error != nil {
			slog.Error(fmt.Sprintf("Error saving transaction for user with id %v: %s", input.UserID, result.Error.Error()))
			return result.Error
		}

		balance = user.Balance

		return nil
	})
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (service *ServiceImpl) TransferMoney(input TransferMoneyInput) error {
	var createdAt time.Time

	invalidTransactionParameters := fmt.Errorf("transaction not possible with sent values")
	var err error
	for attempt := 0; attempt < 10; attempt++ {
		err := service.db.Transaction(func(tx *gorm.DB) error {
			result := tx.Model(&models.User{}).
				Where("id = ? AND balance - ? >= 0", input.FromUserID, input.Amount).
				Update("balance", gorm.Expr("balance - ?", input.Amount))
			if result.Error != nil {
				slog.Error(fmt.Sprintf("Error fetching user with id %v: %s", input.FromUserID, result.Error.Error()))
				return err
			}
			if result.RowsAffected == 0 {
				slog.Error(fmt.Sprintf("User %v not found: %s", input.FromUserID, result.Error.Error()))
				return invalidTransactionParameters
			}

			if err := tx.Model(&models.User{}).
				Where("id = ?", input.ToUserID).
				Update("balance", gorm.Expr("balance + ?", input.Amount)).Error; err != nil {
				slog.Error(fmt.Sprintf("Error fetching user with id %v: %s", input.ToUserID, result.Error.Error()))
				return err
			}

			createdAt = time.Now()

			return nil
		})
		if err == nil || errors.Is(err, invalidTransactionParameters) {
			break
		}

		// Sleep before retrying the transaction
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		return err
	}

	transactions := []models.Transaction{
		{
			UserID: input.FromUserID,
			Amount: -input.Amount,
			Model: gorm.Model{
				CreatedAt: createdAt,
			},
		},
		{
			UserID: input.ToUserID,
			Amount: input.Amount,
			Model: gorm.Model{
				CreatedAt: createdAt,
			},
		},
	}
	if err = service.db.Create(&transactions).Error; err != nil {
		slog.Error(fmt.Sprintf("Error saving transactions: %s", err.Error()))
		return err
	}

	return nil
}
