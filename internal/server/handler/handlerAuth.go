package handler

import (
	"errors"

	"github.com/yury-nazarov/goph_keeper/internal/models"
	"github.com/yury-nazarov/goph_keeper/internal/server/service/auth"
	"net/http"

	"go.uber.org/zap"
)

// AuthController - контроллер обработки HTTP запросов для модуля аутентификаци
type authController struct {
	sessions Sessions
	log      *zap.Logger
	auth     Auth
}

// NewAuthController создает новый экземпляр контроллера который передаем в роутер
func NewAuthController(auth Auth, sessions Sessions, log *zap.Logger) *authController {
	c := &authController{
		auth:     auth,
		sessions: sessions,
		log:      log,
	}
	return c
}


// SignUp регистрация нового пользователя
func (c *authController) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// Извлекаем login, pwd из HTTP запроса
	if err = JSONUnmarshal(r, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Регестрируем нового пользователя, добавляя в структурку поля ID и Token,
	err = c.auth.RegisterUser(r.Context(), &user)
	if errors.Is(err, auth.AuthenticationError) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if errors.Is(err, auth.LoginAlreadyExist) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if errors.Is(err, auth.InternalServerError) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Устанавливаем токен в заголовок и отвечаем клиенту
	w.Header().Set("Authorization", user.Token)
	w.WriteHeader(http.StatusCreated)
}

// SignIn - аутентификация пользователя
func (c *authController) SignIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// Извлекаем login, pwd из HTTP запроса
	if err = JSONUnmarshal(r, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Аутентифицируем пользователя
	err = c.auth.UserLogIn(r.Context(), &user)
	if errors.Is(err, auth.AuthenticationError) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if errors.Is(err, auth.InternalServerError) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", user.Token)
	w.WriteHeader(http.StatusOK)
}

// SignOut - выход пользователя
func (c *authController) SignOut(w http.ResponseWriter, r *http.Request) {
	// Получить токен из заголовка
	token := r.Header.Get("Authorization")
	// Удалить по токену запись в сессиях
	err = c.auth.LogOutUser(r.Context(), token)
	if errors.Is(err, auth.TokenNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, auth.InternalServerError) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Пользователь успешно удален
	w.WriteHeader(http.StatusOK)
}

