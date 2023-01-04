package postgres

import (
	"context"
	"database/sql"
	"github.com/yury-nazarov/goph_keeper/internal/app/models"
	"github.com/yury-nazarov/goph_keeper/internal/options"
	"strconv"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose"
	"go.uber.org/zap"
)

// Реализация создания нового подклчения к БД

// DB интерфейс для работы с Psql
type DB interface {
	UserExist(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, login string, password string) (int, error)
	UserIsValid(ctx context.Context, user models.User) (int, error)

	AddSecret(ctx context.Context, secret models.Secret) (int, error)
	GetSecretList(ctx context.Context, userID int) ([]models.Secret, error)
	GetSecretByID(ctx context.Context, secret models.Secret) (models.Secret, error)
	UpdateSecretByID(ctx context.Context, secret models.Secret) error
	DeleteSecretByID(ctx context.Context, secret models.Secret) error

	Close() error
}

// psql описывает поля необходимые для работы с БД
type psql struct {
	db  *sql.DB
	log *zap.Logger
}

// NewPostgres  инициирует подключение к БД
func NewPostgres(log *zap.Logger, cfg options.Config) (*psql, error) {
	// Открываем подключение к БД
	db, err := sql.Open("pgx", cfg.DB)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение к БД
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Формируем объект для работы с БД
	p := &psql{db: db, log: log}

	// Запускаем миграции
	err = p.migrations(cfg.MigrateTo, cfg.MigrateFile)
	if err != nil {
		return nil, err
	}

	log.Info("DB storage connection success")
	return p, nil
}

// migrations запускает миграции
//
func (p *psql) migrations(migrateTo string, migrateFile string) error {
	var err error
	var migrationDownToNum int64

	if len(migrateTo) > 0 {
		migrationDownToNum, err = strconv.ParseInt(migrateTo, 10, 64)
		if err != nil {
			return err
		}

		err = goose.DownTo(p.db, migrateFile, migrationDownToNum)
	} else {
		err = goose.Up(p.db, migrateFile)
	}

	if err != nil {
		return err
	}
	return nil
}

