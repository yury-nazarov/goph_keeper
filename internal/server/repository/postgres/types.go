package postgres

import (
	"context"

	"github.com/yury-nazarov/goph_keeper/internal/models"
)

// DB импортируемый интерфейс для клозера
// 	  Все поля нужны для запуска тестов repository слоя
type DB interface {
	UserExist(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, login string, password string) (int, error)
	UserIsValid(ctx context.Context, user models.User) (int, error)

	AddSecret(ctx context.Context, secret models.Secret) (int, error)
	GetSecretList(ctx context.Context, userID int) ([]models.Secret, error)
	GetSecretByID(ctx context.Context, secret models.Secret) (models.Secret, error)
	UpdateSecretByID(ctx context.Context, secret models.Secret) error
	DeleteSecretByID(ctx context.Context, secret models.Secret) error

	Close() error
}
