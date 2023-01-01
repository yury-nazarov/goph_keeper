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
	UserLogIn(ctx context.Context, user *models.User) error
	LogOutUser(ctx context.Context, token string) error
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
	// Проверяем наличие логина в БД
	ok, err := a.db.UserExist(ctx, user.Login)
	if err != nil {
		a.log.Warn("Сan`t check user exist",
			zap.String("method", "Auth.RegisterUser"),
			zap.String("error", err.Error()),
			zap.String("user.Login", user.Login))
		return tools.NewErr500("")
	}
	if ok {
		a.log.Info("Сan`t register new user",
			zap.String("method", "Auth.RegisterUser"),
			zap.String("error", "username already exists"),
			zap.String("user.Login", user.Login))
		return tools.NewErr409("")
	}

	// Создаем пользователя
	user.Password = a.hashPassword(user.Password)
	user.ID, err = a.db.CreateUser(ctx, user.Login, user.Password)
	if err != nil {
		a.log.Warn("Сan`t create new user",
			zap.String("method", "Auth.RegisterUser"),
			zap.String("error", err.Error()),
			zap.String("user.Login", user.Login))
		return tools.NewErr500("")
	}

	// Создаем токен
	user.Token, err = a.createToken()
	if err != nil {
		a.log.Warn("Сan`t create token",
			zap.String("method", "Auth.RegisterUser"),
			zap.String("error", err.Error()),
			zap.String("user.Login", user.Login))
		return tools.NewErr500("")
	}

	// Логинем пользователя
	err = a.sessions.AddToken(ctx, user.Token, user.ID)
	if err != nil {
		a.log.Warn("Сan`t add token to session",
			zap.String("method", "Auth.RegisterUser"),
			zap.String("error", err.Error()),
			zap.String("user.Login", user.Login))
		return tools.NewErr500("")
	}
	// Успешное создане пользователя
	a.log.Info("Success registered new user",
		zap.String("method", "Auth.RegisterUser"),
		zap.String("user.Login", user.Login),
		zap.Int("user.ID", user.ID),
		zap.String("user.Password[hash]", user.Password),
		zap.String("user.Token", user.Token))
	return nil
}

// UserLogIn - описывает процедуру входа пользователя
func (a *auth) UserLogIn(ctx context.Context, user *models.User) error {
	var err error

	// Хешируем пароль
	user.Password = a.hashPassword(user.Password)
	a.log.Debug("Success create hashPwd",
		zap.String("method", "Auth.UserLogIn"),
		zap.String("user.Login", user.Login),
		zap.String("hashPassword", user.Password))

	// Проверяем совпадают ли логин/пароль для пользователя с теми что в БД
	user.ID, err = a.db.UserIsValid(ctx, *user)
	if err != nil {
		a.log.Info("Incorrect username or password",
			zap.String("method", "Auth.UserLogIn"),
			zap.String("user.Login", user.Login),
			zap.String("error", err.Error()))
		return tools.NewErr401("")
	}

	// Создаем токен
	user.Token, err = a.createToken()
	if err != nil {
		a.log.Warn("Сan`t create token",
			zap.String("method", "Auth.UserLogIn"),
			zap.String("user.Login", user.Login),
			zap.String("error", err.Error()))
		return tools.NewErr500("")
	}

	// Логинем пользователя
	err = a.sessions.AddToken(ctx, user.Token, user.ID)
	if err != nil {
		a.log.Warn("Can't add token to session",
			zap.String("method", "Auth.UserLogIn"),
			zap.String("user.Login", user.Login),
			zap.String("error", err.Error()))
		return tools.NewErr500("")
	}
	// Успешный логин
	a.log.Info("Success LogIn to system",
		zap.String("method", "Auth.UserLogIn"),
		zap.String("user.Login", user.Login),
		zap.Int("user.ID", user.ID),
		zap.String("user.Password[hash]", user.Password),
		zap.String("user.Token", user.Token))
	return nil
}

// LogOutUser логика выхода пользователя и удаления токена в сессиях
func (a *auth) LogOutUser(ctx context.Context, token string) error {
	// Проверяем наличие залогиненой сессии
	userID, err := a.sessions.GetUserID(ctx, token)
	if err != nil {
		a.log.Info("Incorrect token",
			zap.String("method", "Auth.LogOutUser"),
			zap.String("token", token),
			zap.String("error", err.Error()))
		return tools.NewErr404( "")
	}

	// Если есть, удаляем сессию
	err = a.sessions.DeleteToken(ctx, token)
	if err != nil {
		a.log.Warn("Can`t delete token",
			zap.String("method", "Auth.LogOutUser"),
			zap.Int("user.ID", userID),
			zap.String("user.Token", token),
			zap.String("error", err.Error()))
		return tools.NewErr500("")
	}
	// Успешный логаут
	a.log.Info("Success logOut",
		zap.String("method", "Auth.LogOutUser"),
		zap.Int("user.ID", userID),
		zap.String("user.Token", token))
	return nil
}

// CreateUserToken создание пользовательского токена
func (a *auth) createToken() (string, error) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		msg := fmt.Sprintf("Can`t create token: %s", err.Error())
		return "", fmt.Errorf(msg)
	}
	// Успешное создание токена
	token := fmt.Sprintf("%x", b)
	return token, nil

}

// hashPassword хеш из пароля
func (a *auth) hashPassword(password string) string {
	hashPwd := md5.Sum([]byte(password))
	hp := fmt.Sprintf("%x", hashPwd)
	return hp
}
