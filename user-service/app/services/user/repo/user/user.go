package user

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"user-service/models"
)

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *RepositoryImpl {
	return &RepositoryImpl{
		db: db,
	}
}

func (repo *RepositoryImpl) Insert(email string) (models.User, error) {
	var user models.User
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		user = models.User{ID: uuid.New(), Email: email}
		result := repo.db.Create(&user)
		if result.Error != nil {
			slog.Error(fmt.Sprintf("Failed to create user with email %s: %s", email, result.Error.Error()))
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
			slog.Error(fmt.Sprintf("Failed to delete user with id %v: %s", id, result.Error.Error()))
			return result.Error
		}

		return nil
	})
}

func (repo *RepositoryImpl) GetByEmail(email string) (models.User, error) {
	var user models.User
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		result := repo.db.Clauses(
			clause.Locking{Strength: clause.LockingStrengthShare},
		).Where("email = ?", email).First(&user)
		if result.Error != nil {
			slog.Error(fmt.Sprintf("Failed to get user with email %s: %s", email, result.Error.Error()))
			return result.Error
		}

		return nil
	})
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
