package auth

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"fmt"

	"github.com/yury-nazarov/goph_keeper/internal/app/models"
	"github.com/yury-nazarov/goph_keeper/internal/app/repository"
	"github.com/yury-nazarov/goph_keeper/pkg/tools"

	"go.uber.org/zap"
)

// Модуль для работы с авторизацией и аутентификацией пользователей

type Auth interface {
	RegisterUser(ctx context.Context, user *models.User) error
}

// auth структурка для работы модуля аутентификации/авторизации
type auth struct {
	log      *zap.Logger
	sessions repository.Sessions
	db       repository.DB
}

// New создает объект на основе стурктуры auth
func New(log *zap.Logger, sessions repository.Sessions, db repository.DB) *auth {
	a := &auth{
		log:      log,
		sessions: sessions,
		db:       db,
	}
	return a
}

// RegisterUser - логика создания нового пользователя
func (a *auth) RegisterUser(ctx context.Context, user *models.User) error {
	var msg string
	// Проверяем наличие логина в БД
	ok, err := a.db.UserExist(ctx, user.Login)
	if err != nil {
		msg = fmt.Sprintf("can't check user exist: %s", err.Error())
		a.log.Warn(msg)
		return tools.NewErr500(msg)
	}
	if ok {
		msg = fmt.Sprintf("%s - user exist", user.Login)
		a.log.Info(msg)
		return tools.NewErr409(msg)
	}

	// Создаем пользователя
	user.Password = hashPassword(user.Password)
	user.ID, err = a.db.CreateUser(ctx, user.Login, user.Password)
	if err != nil {
		msg = fmt.Sprintf("can't create user: %s", err.Error())
		a.log.Warn(msg)
		return tools.NewErr500(msg)
	}

	// Создаем токен
	user.Token, err = a.createToken()
	if err != nil {
		msg = fmt.Sprintf("can't create token: %s", err.Error())
		a.log.Warn(msg)
		return tools.NewErr500(msg)
	}

	// Логинем пользователя
	err = a.sessions.AddToken(ctx, user.Token, user.ID)
	if err != nil {
		msg = fmt.Sprintf("can't add token to session: %s", err.Error())
		a.log.Warn(msg)
		return tools.NewErr500(msg)
	}
	return nil
}

// CreateUserToken создание пользовательского токена
func (a *auth) createToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		msg := fmt.Sprintf("can't create token: %s", err.Error())
		return "", fmt.Errorf(msg)
	}
	return fmt.Sprintf("%x", b), nil

}

// hashPassword хеш из пароля
func hashPassword(password string) string {
	hashPwd := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hashPwd)
}
