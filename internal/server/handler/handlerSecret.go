package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yury-nazarov/goph_keeper/internal/models"
	"github.com/yury-nazarov/goph_keeper/internal/server/repository/options"
	"github.com/yury-nazarov/goph_keeper/internal/server/service/secret"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// secret - контроллер обработки HTTP запросов
type secretController struct {
	cgf      options.Config
	log      *zap.Logger
	secret   Secret
}

// NewSecretController создает новый экземпляр контроллера который передаем в роутер
func NewSecretController(secret Secret, log *zap.Logger) *secretController {
	c := &secretController{
		secret:   secret,
		log:      log,
	}
	return c
}

// SecretNew - создание нового секрета
func (c *secretController) SecretNew(w http.ResponseWriter, r *http.Request) {
	var item models.Secret
	// Извлекаем name, data, description из HTTP запроса
	if err = JSONUnmarshal(r, &item); err != nil {
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
func (c *secretController) SecretList(w http.ResponseWriter, r *http.Request) {
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
func (c *secretController) GetSecretByID(w http.ResponseWriter, r *http.Request) {
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
func (c *secretController) UpdateSecretByID(w http.ResponseWriter, r *http.Request) {
	var secretItem models.Secret
	// Получаем userID который добавляем в контекст из middleware
	secretItem.UserID = r.Context().Value("userID").(int)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Извлекаем id, name, data, description из HTTP запроса
	if err = JSONUnmarshal(r, &secretItem); err != nil {
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
func (c *secretController) DeleteSecretByID(w http.ResponseWriter, r *http.Request) {
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
