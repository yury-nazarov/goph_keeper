package secret

import (
	"context"
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

