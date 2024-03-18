package user

import (
	"github.com/google/uuid"
	"user-service/models"
)

type Repository interface {
	Insert(email string) (models.User, error)
	Delete(id uuid.UUID) error
	GetByEmail(email string) (models.User, error)
}
