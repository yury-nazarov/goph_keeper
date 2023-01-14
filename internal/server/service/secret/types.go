package secret

import (
	"context"

	"github.com/yury-nazarov/goph_keeper/internal/models"
)

// DB интерфейс для работы с Psql
type DB interface {
	AddSecret(ctx context.Context, secret models.Secret) (int, error)
	GetSecretList(ctx context.Context, userID int) ([]models.Secret, error)
	GetSecretByID(ctx context.Context, secret models.Secret) (models.Secret, error)
	UpdateSecretByID(ctx context.Context, secret models.Secret) error
	DeleteSecretByID(ctx context.Context, secret models.Secret) error
}
