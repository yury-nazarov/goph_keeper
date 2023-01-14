package handler

import (
	"context"
	"github.com/yury-nazarov/goph_keeper/internal/models"
)

// Работа с репозиторием

// Sessions - интерфейс для работы с сессиями
type Sessions interface {
	AddToken(ctx context.Context, token string, userID int) error
	GetUserID(ctx context.Context, token string) (int, error)
	DeleteToken(ctx context.Context, token string) error
}

// Работа с бизнес логикой

// Auth аутентификация пользователя
type Auth interface {
	RegisterUser(ctx context.Context, user *models.User) error
	UserLogIn(ctx context.Context, user *models.User) error
	LogOutUser(ctx context.Context, token string) error
}

type Secret interface {
	Create(ctx context.Context, secret models.Secret) error
	List(ctx context.Context, userID int) ([]models.Secret, error)
	GetByID(ctx context.Context, secretID int) (models.Secret, error)
	PutByID(ctx context.Context, item models.Secret) error
	DeleteByID(ctx context.Context, secretID int) error
}