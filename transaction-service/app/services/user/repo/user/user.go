package user

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	"transaction-service/models"
)

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *RepositoryImpl {
	return &RepositoryImpl{
		db: db,
	}
}

func (repo *RepositoryImpl) Insert(id uuid.UUID, createdAt time.Time) (models.User, error) {
	user := models.User{ID: id, Balance: 0.0, CreatedAt: createdAt}
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		result := repo.db.Create(&user)
		if result.Error != nil {
			slog.Error(fmt.Sprintf("Error creating user with id %v: %v", id, result.Error.Error()))
			return result.Error
		}

		return nil
	})
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repo *RepositoryImpl) Delete(id uuid.UUID) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		result := repo.db.Delete(&models.User{}, id)
		if result.Error != nil {
			slog.Error(fmt.Sprintf("Error deleting user with id %v: %v", id, result.Error.Error()))
			return result.Error
		}

		return nil
	})
}

func (repo *RepositoryImpl) Get(id string) (models.User, error) {
	var user models.User
	result := repo.db.Clauses(clause.Locking{
		Strength: clause.LockingStrengthShare,
	}).Where("id = ?", id).First(&user)
	if result.Error != nil {
		slog.Error(fmt.Sprintf("Failed to get user with id %v: %s", id, result.Error.Error()))
		return models.User{}, result.Error
	}

	return user, nil
}
