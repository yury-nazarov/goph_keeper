package handler

import (
	"errors"
	"net/http"

	"github.com/yury-nazarov/goph_keeper/internal/app/models"
	"github.com/yury-nazarov/goph_keeper/internal/app/repository"
	"github.com/yury-nazarov/goph_keeper/internal/app/service/auth"
	"github.com/yury-nazarov/goph_keeper/internal/options"
	"github.com/yury-nazarov/goph_keeper/pkg/tools"

	"go.uber.org/zap"
)

// Общие переменные для хендлеров
var (
	err    error
	err409 *tools.Err409
	err500 *tools.Err500
)

// Controller - контроллер обработки HTTP запросов
type Controller struct {
	db       repository.DB
	sessions repository.Sessions
	cgf      options.Config
	log      *zap.Logger
	auth     auth.Auth
}

// NewController создает новый экземпляр контроллера который передаем в роутер
func NewController(db repository.DB, sessions repository.Sessions, cfg options.Config, log *zap.Logger, auth auth.Auth) *Controller {
	c := &Controller{
		db:       db,
		sessions: sessions,
		cgf:      cfg,
		log:      log,
		auth:     auth,
	}
	return c
}

func (c *Controller) Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("v001"))
	return
}

// SignUp регистрация нового пользователя
func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// Извлекаем login, pwd из HTTP запроса
	if err = tools.JSONUnmarshal(r, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Регестрируем нового пользователя, добавляя в структурку поля ID и Token,
	err = c.auth.RegisterUser(r.Context(), &user)
	if errors.As(err, &err409) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if errors.As(err, &err500) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Устанавливаем токен в заголовок и отвечаем клиенту
	w.Header().Set("Authorization", user.Token)
	w.WriteHeader(http.StatusCreated)
}

// SignIn - аутентификация пользователя
func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// Извлекаем login, pwd из HTTP запроса
	if err = tools.JSONUnmarshal(r, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Аутентифицируем пользователя
	err = c.auth.UserLogIn(r.Context(), &user)
	if err != nil {
		//
	}
}
