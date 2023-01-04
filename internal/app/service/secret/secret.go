package secret

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/yury-nazarov/goph_keeper/internal/app/models"
	"github.com/yury-nazarov/goph_keeper/internal/app/repository/postgres"
	"github.com/yury-nazarov/goph_keeper/pkg/tools"
	"go.uber.org/zap"
)

var err error

type secret struct {
	db  postgres.DB
	log *zap.Logger
}

type Secret interface {
	Create(ctx context.Context, secret models.Secret) error
	List(ctx context.Context, userID int) ([]models.Secret, error)
	GetByID(ctx context.Context, secretID int) (models.Secret, error)
	PutByID(ctx context.Context, item models.Secret) error
	DeleteByID(ctx context.Context, secretID int) error
}

func NewSecret(db postgres.DB, logger *zap.Logger) *secret {
	s := &secret{
		db: db,
		log: logger,
	}
	return s
}

// Create логика создания нового секрета
func (s *secret) Create(ctx context.Context, secret models.Secret) error {
	secret.ID, err = s.db.AddSecret(ctx, secret)
	if err != nil {
		s.log.Warn("Can't add secret",
			zap.String("method", "Secret.Create"),
			zap.Int("user.ID", secret.UserID),
			zap.String("error", err.Error()))
		return tools.NewErr500("")
	}
	s.log.Info("Success create secret",
		zap.String("method", "Secret.Create"),
		zap.Int("secret.UserID", secret.UserID),
		zap.Int("secret.ID", secret.ID),
		zap.String("secret.Name", secret.Name))
	return nil
}

// List логика получения всех секретов пользователя
func (s *secret) List(ctx context.Context, userID int) (secrets []models.Secret, err error) {
	secrets, err = s.db.GetSecretList(ctx, userID)
	if err != nil {
		s.log.Warn("Can't get list of secrets",
			zap.String("method", "Secret.List"),
			zap.Int("userID", userID),
			zap.String("error", err.Error()))
		return nil, tools.NewErr500("")
	}
	s.log.Info("Success get list of secret",
		zap.String("method", "Secret.List"),
		zap.Int("userID", userID))
	s.log.Debug("Success get list of secret",
		zap.String("method", "Secret.List"),
		zap.Int("userID", userID),
		zap.String("secrets", fmt.Sprintf("%+v", secrets)))
	return secrets, nil
}

// GetByID вернет секрет по ID
func (s *secret) GetByID(ctx context.Context, secretID int) (secret models.Secret, err error) {
	// Получаем userID который добавляем в контекст из middleware
	secret.UserID = ctx.Value("userID").(int)
	s.log.Debug("secret.GetByID", zap.Int("userID", secret.UserID))

	secret.ID = secretID

	secret, err = s.db.GetSecretByID(ctx, secret)
	if errors.Is(err, sql.ErrNoRows) {
		s.log.Warn("can't get secret",
			zap.String("method", "secret.GetByID"),
			zap.Int("userID", secret.UserID),
			zap.Int("secretID", secret.ID),
			zap.String("error", err.Error()),
		)
		return secret, tools.NewErr404("")
	}
	if err != nil {
		s.log.Warn("can't get secret",
			zap.String("method", "secret.GetByID"),
			zap.Int("userID", secret.UserID),
			zap.Int("secretID", secret.ID),
			zap.String("error", err.Error()),
			)
		return secret, tools.NewErr500("")
	}
	return secret, nil
}

func (s *secret) PutByID(ctx context.Context, item models.Secret) error {
	var item2 models.Secret
	// Проверяем что секрет пренадлежит этомй пользователю
	item2, err = s.db.GetSecretByID(ctx, item)
	if err  != nil || item.ID != item2.ID || item.UserID != item2.UserID{
		s.log.Warn("HTTP request isn`t authorized",
			zap.String("method", "secret.PutByID"),
			zap.Int("userID", item.UserID),
			zap.Int("secretID", item.ID),
			zap.String("error", err.Error()))
		return tools.NewErr401("")
	}

	// Обновляем серет в БД
	err = s.db.UpdateSecretByID(ctx, item)
	if err != nil {
		s.log.Warn("can't update secret",
			zap.String("method", "secret.PutByID"),
			zap.Int("userID", item.UserID),
			zap.Int("secretID", item.ID),
			zap.String("error", err.Error()),
		)
		return tools.NewErr500("")
	}
	return nil
}

// DeleteByID удаляет секрет
func (s *secret) DeleteByID(ctx context.Context, secretID int) error {
	var (
		item 	models.Secret
		item2 	models.Secret
	)
	// Получаем userID который добавляем в контекст из middleware
	item.UserID = ctx.Value("userID").(int)
	item.ID = secretID

	// Проверяем что секрет пренадлежит этомй пользователю
	item2, err = s.db.GetSecretByID(ctx, item)
	if err  != nil || item.ID != item2.ID || item.UserID != item2.UserID{
		s.log.Warn("HTTP request isn`t authorized",
			zap.String("method", "secret.PutByID"),
			zap.Int("userID", item.UserID),
			zap.Int("secretID", item.ID),
			zap.String("error", err.Error()))
		return tools.NewErr401("")
	}

	// Удаляем запись в БД
	err = s.db.DeleteSecretByID(ctx, item)
	if err != nil {
		s.log.Warn("can't delete secret",
			zap.String("method", "secret.PutByID"),
			zap.Int("userID", item.UserID),
			zap.Int("secretID", item.ID),
			zap.String("error", err.Error()))
		return tools.NewErr500("")
	}
	return nil
}