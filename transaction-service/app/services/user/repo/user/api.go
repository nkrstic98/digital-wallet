package user

import (
	"github.com/google/uuid"
	"time"
	"transaction-service/models"
)

type Repository interface {
	Insert(id uuid.UUID, createdAt time.Time) (models.User, error)
	Delete(id uuid.UUID) error
	Get(id string) (models.User, error)
}
