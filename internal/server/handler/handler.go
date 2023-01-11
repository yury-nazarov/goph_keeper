package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yury-nazarov/goph_keeper/internal/models"
	"github.com/yury-nazarov/goph_keeper/internal/server/repository/inmemory"
	"github.com/yury-nazarov/goph_keeper/internal/server/repository/options"
	"github.com/yury-nazarov/goph_keeper/internal/server/repository/postgres"
	"github.com/yury-nazarov/goph_keeper/internal/server/service/auth"
	"github.com/yury-nazarov/goph_keeper/internal/server/service/secret"
	"github.com/yury-nazarov/goph_keeper/pkg/tools"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var err    error

// Controller - контроллер обработки HTTP запросов
type Controller struct {
	db       postgres.DB
	sessions inmemory.Sessions
	cgf      options.Config
	log      *zap.Logger
	auth     auth.Auth
	secret   secret.Secret
}

// NewController создает новый экземпляр контроллера который передаем в роутер
func NewController(db postgres.DB, sessions inmemory.Sessions, cfg options.Config, log *zap.Logger, auth auth.Auth, secret secret.Secret) *Controller {
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
func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// Извлекаем login, pwd из HTTP запроса
	if err = tools.JSONUnmarshal(r, &user); err != nil {
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
func (c *Controller) SignOut(w http.ResponseWriter, r *http.Request) {
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

// SecretNew - создание нового секрета
func (c *Controller) SecretNew(w http.ResponseWriter, r *http.Request) {
	var item models.Secret
	// Извлекаем name, data, description из HTTP запроса
	if err = tools.JSONUnmarshal(r, &item); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Получаем userID который добавляем в контекст из middleware
	item.UserID = r.Context().Value("userID").(int)
	c.log.Debug("handler.SecretNew", zap.String("struct debug", fmt.Sprintf("%+v", item)))

	// Оправляем в слой бизнес логики для создания секрета в БД
	err = c.secret.Create(r.Context(), item)
	if errors.Is(err, secret.InternalServerError) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// SecretList - получить список секретов
func (c *Controller) SecretList(w http.ResponseWriter, r *http.Request) {
	var secrets []models.Secret
	var secretsJSON []byte

	// Получаем userID который добавляем в контекст из middleware
	userID := r.Context().Value("userID").(int)
	c.log.Debug("handler.SecretList", zap.Int("userID", userID))

	secrets, err = c.secret.List(r.Context(), userID)
	if errors.Is(err, secret.InternalServerError) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Замаршалить в JSON
	secretsJSON, err = json.Marshal(secrets)
	if err != nil {
		c.log.Warn("can't marshal to JSON",
			zap.String("method", "handler.SecretList"),
			zap.Int("userID", userID),
			zap.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Список секретов успешно отправлен пользователю
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(secretsJSON)
	if err != nil {
		c.log.Warn("can't write response to client",
			zap.String("method", "handler.SecretList"),
			zap.Int("userID", userID),
			zap.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetSecretByID вернет секрет по ID
func (c *Controller) GetSecretByID(w http.ResponseWriter, r *http.Request) {
	var (
		item       models.Secret
		secretJSON []byte
		secretID   int
	)

	// Получаем secretID из URL
	secretID, err = strconv.Atoi(chi.URLParam(r, "secretID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Получаем секрет
	item, err = c.secret.GetByID(r.Context(), secretID)
	if errors.Is(err, secret.ItemNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, secret.InternalServerError) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Замаршалить в JSON
	secretJSON, err = json.Marshal(item)
	if err != nil {
		c.log.Warn("can't marshal to JSON",
			zap.String("method", "handler.SecretByID"),
			zap.Int("userID", item.UserID),
			zap.Int("secretID", secretID),
			zap.String("secret", fmt.Sprintf("%+v", item)),
			zap.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Секрет успешно отправлен пользователю
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(secretJSON)
	if err != nil {
		c.log.Warn("can't write response to client",
			zap.String("method", "handler.SecretByID"),
			zap.Int("userID", item.UserID),
			zap.Int("secretID", secretID),
			zap.String("secret", fmt.Sprintf("%+v", item)),
			zap.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UpdateSecretByID обновление секрета по ID
func (c *Controller) UpdateSecretByID(w http.ResponseWriter, r *http.Request) {
	var secretItem models.Secret
	// Получаем userID который добавляем в контекст из middleware
	secretItem.UserID = r.Context().Value("userID").(int)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Извлекаем id, name, data, description из HTTP запроса
	if err = tools.JSONUnmarshal(r, &secretItem); err != nil {
		c.log.Warn("can't unmarshal json",
			zap.String("method", "handler.UpdateSecretByID"),
			zap.Int("userID", secretItem.UserID),
			zap.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отправляем в слой бизнесл логики для доп. проверок и обновления
	err = c.secret.PutByID(r.Context(), secretItem)
	if errors.Is(err, secret.AuthenticationError) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if errors.Is(err, secret.ItemNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if errors.Is(err, secret.InternalServerError) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// DeleteSecretByID удалить секрет
func (c *Controller) DeleteSecretByID(w http.ResponseWriter, r *http.Request) {
	var secretID int
	// Получаем secretID из URL
	secretID, err = strconv.Atoi(chi.URLParam(r, "secretID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отправляем в слой бизнесл логики для доп. проверок и обновления
	err = c.secret.DeleteByID(r.Context(), secretID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
