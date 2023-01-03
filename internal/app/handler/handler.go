package handler

import (
	"errors"
	"fmt"
	"github.com/yury-nazarov/goph_keeper/internal/app/models"
	"github.com/yury-nazarov/goph_keeper/internal/app/repository"
	"github.com/yury-nazarov/goph_keeper/internal/app/service/auth"
	"github.com/yury-nazarov/goph_keeper/internal/app/service/secret"
	"github.com/yury-nazarov/goph_keeper/internal/options"
	"github.com/yury-nazarov/goph_keeper/pkg/tools"
	"net/http"

	"go.uber.org/zap"
)

// Общие переменные для хендлеров
var (
	err    error
	err401 *tools.Err401
	err404 *tools.Err404
	err409 *tools.Err409
	err500 *tools.Err500
)

// Controller - контроллер обработки HTTP запросов
type Controller struct {
	db       	repository.DB
	sessions 	repository.Sessions
	cgf      	options.Config
	log      	*zap.Logger
	auth     	auth.Auth
	secret     	secret.Secret
}

// NewController создает новый экземпляр контроллера который передаем в роутер
func NewController(db repository.DB, sessions repository.Sessions, cfg options.Config, log *zap.Logger, auth auth.Auth, secret secret.Secret) *Controller {
	c := &Controller{
		db:       db,
		sessions: sessions,
		cgf:      cfg,
		log:      log,
		auth:     auth,
		secret:   secret,
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
	if errors.As(err, &err401) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if errors.As(err, &err500) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", user.Token)
	w.WriteHeader(http.StatusOK)
}


// SignOut - выход пользователя
func (c *Controller) SignOut(w http.ResponseWriter, r *http.Request) {
	// Получить токен из заголовка
	token := r.Header.Get("Authorization")
	// Удалить по токену запись в сессиях
	err = c.auth.LogOutUser(r.Context(), token)
	if errors.As(err, &err404) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.As(err, &err500) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Пользователь успешно удален
	w.WriteHeader(http.StatusOK)
}

// SecretNew - создание нового секрета
func (c *Controller) SecretNew(w http.ResponseWriter, r *http.Request) {
	var secret models.Secret
	// Извлекаем name, data, description из HTTP запроса
	if err = tools.JSONUnmarshal(r, &secret); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Добавляем userID который добавляем в контекст из middleware
	secret.UserID = r.Context().Value("userID").(int)
	c.log.Debug("handler.SecretNew", zap.String("struct debug", fmt.Sprintf("%+v", secret)))

	// Оправляем в слой бизнес логики для создания секрета в БД
	err = c.secret.Create(r.Context(), secret)
	if errors.As(err, &err500) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}