package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/yury-nazarov/goph_keeper/internal/models"
	"github.com/yury-nazarov/goph_keeper/internal/server/repository/options"

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

// New  инициирует подключение к БД
func New(log *zap.Logger, cfg options.Config) (*psql, error) {
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
func (p *psql) migrations(migrateTo string, migrateFile string) error {
	var err error
	var migrateToNum int64

	if len(migrateTo) > 0 {
		migrateToNum, err = strconv.ParseInt(migrateTo, 10, 64)
		if err != nil {
			return err
		}

		err = goose.DownTo(p.db, migrateFile, migrateToNum)
	} else {
		err = goose.Up(p.db, migrateFile)
	}

	if err != nil {
		return err
	}
	return nil
}



// UserExist проверяет наличие пользвоателя в БД, вернет:
//			 true, nil -  если пользователь уже есть в БД
//           false, nil - если пользователя нет в БД
//			 false, err - если произошла ошибка во время выполнения запроса
func (p *psql) UserExist(ctx context.Context, login string) (bool, error) {
	var loginFromDB string
	err := p.db.QueryRowContext(ctx, `SELECT login FROM app_user WHERE login=$1 LIMIT 1`, login).Scan(&loginFromDB)
	// Записи нет в БД
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if len(login) > 0 {
		return true, nil
	}
	return false, nil
}

// CreateUser создает нового пользователя
func (p *psql) CreateUser(ctx context.Context, login string, password string) (int, error) {
	var userID int
	err := p.db.QueryRowContext(ctx,
		`INSERT INTO app_user (login, password) VALUES ($1, $2) RETURNING id`, login, password).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf(err.Error())
	}
	return userID, nil
}

// UserIsValid проверяет валидный логин пароль для пользователя
// 				вернет nil - если УЗ принаделжит пользователю
func (p *psql) UserIsValid(ctx context.Context, user models.User) (int, error) {
	err := p.db.QueryRowContext(ctx,
		`SELECT id FROM app_user WHERE login=$1 and password=$2 LIMIT 1`, user.Login, user.Password).Scan(&user.ID)
	// Записи нет в БД
	if errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	if err != nil {
		return 0, err
	}
	// если запись найдена по логину и хешу пароля, то считаем, что учетные данные валидны
	return user.ID, nil
}

// AddSecret создает в БД запись для нового секрета
func (p *psql) AddSecret(ctx context.Context, secret models.Secret) (int, error) {
	err := p.db.QueryRowContext(ctx, `INSERT INTO app_secret (user_id, name, data, description) VALUES ($1, $2, $3, $4) RETURNING id`, secret.UserID, secret.Name, secret.Data, secret.Description).Scan(&secret.ID)
	if err != nil {
		return 0, err
	}
	return secret.ID, nil
}

// GetSecretList получает из БД список секретов по userID
func (p *psql) GetSecretList(ctx context.Context, userID int) (secretList []models.Secret, err error)  {
	var secret models.Secret
	rows, err := p.db.QueryContext(ctx, `SELECT id, user_id, name, data, description FROM app_secret WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&secret.ID, &secret.UserID, &secret.Name, &secret.Data, &secret.Description); err != nil {
			p.log.Warn("can't read string from query",
				zap.String("method", "psql.GetSecretList"),
				zap.Int("userID", userID),
				zap.String("error", err.Error()))
		}
		secretList = append(secretList, secret)
	}
	return secretList, nil
}

// GetSecretByID получает из БД секрет по ID для конкретного пользователя
func (p *psql) GetSecretByID(ctx context.Context, secret models.Secret) (models.Secret, error) {
	err := p.db.QueryRowContext(ctx, `SELECT name, data, description FROM app_secret WHERE id=$1 AND user_id=$2 LIMIT 1`, secret.ID, secret.UserID).Scan(&secret.Name, &secret.Data, &secret.Description)
	// Записи нет в БД
	if errors.Is(err, sql.ErrNoRows) {
		return secret, err
	}
	if err != nil {
		return secret, err
	}
	return secret, nil
}

// UpdateSecretByID обновляет секрет
func (p *psql) UpdateSecretByID(ctx context.Context, secret models.Secret) error {
	_, err := p.db.ExecContext(ctx, `UPDATE app_secret SET name=$1, data=$2, description=$3 WHERE id=$4 AND user_id=$5`, secret.Name, secret.Data, secret.Description, secret.ID, secret.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (p *psql) DeleteSecretByID(ctx context.Context, secret models.Secret) error {
	_, err := p.db.ExecContext(ctx, `DELETE FROM app_secret WHERE id=$1 AND user_id=$2`, secret.ID, secret.UserID)
	if err != nil {
		return err
	}
	return nil
}

// Close закрываем соединение к БД
func (p *psql) Close() error {
	err := p.db.Close()
	if err != nil {
		return fmt.Errorf("can't close DB connection. Error: %s", err.Error())
	}
	return err
}
