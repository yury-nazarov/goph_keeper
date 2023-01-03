package secret

import (
	"context"
	"fmt"
	"github.com/yury-nazarov/goph_keeper/internal/app/models"
	"github.com/yury-nazarov/goph_keeper/internal/app/repository"
	"github.com/yury-nazarov/goph_keeper/pkg/tools"
	"go.uber.org/zap"
)

var err error

type secret struct {
	db repository.DB
	log *zap.Logger
}

type Secret interface {
	Create(ctx context.Context, secret models.Secret) error
	List(ctx context.Context, userID int) ([]models.Secret, error)
}

func NewSecret(db repository.DB, logger *zap.Logger) *secret {
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
