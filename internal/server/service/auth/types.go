package auth

import (
	"context"

	"github.com/yury-nazarov/goph_keeper/internal/models"
)

// DB интерфейс для работы с Psql
type DB interface {
	UserExist(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, login string, password string) (int, error)
	UserIsValid(ctx context.Context, user models.User) (int, error)
}


// Sessions - интерфейс для работы с кешем
type Sessions interface {
	AddToken(ctx context.Context, token string, userID int) error
	GetUserID(ctx context.Context, token string) (int, error)
	DeleteToken(ctx context.Context, token string) error
}