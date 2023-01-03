package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/yury-nazarov/goph_keeper/internal/app/models"

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

	Close() error
}

// psql описывает поля необходимые для работы с БД
type psql struct {
	db  *sql.DB
	log *zap.Logger
}

// NewPostgres  инициирует подключение к БД
func NewPostgres(log *zap.Logger, connString string) (*psql, error) {
	// Открываем подключение к БД
	db, err := sql.Open("pgx", connString)
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
	err = p.migrations()
	if err != nil {
		return nil, err
	}

	log.Info("DB storage connection success")
	return p, nil
}

// migrations запускает миграции
func (p *psql) migrations() error {
	err := goose.Up(p.db, "./internal/migrations")
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
	rows, err := p.db.QueryContext(ctx, `SELECT id, name, data, description FROM app_secret WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&secret.ID, &secret.Name, &secret.Data, &secret.Description); err != nil {
			p.log.Warn("can't read string from query",
				zap.String("method", "psql.GetSecretList"),
				zap.Int("userID", userID),
				zap.String("error", err.Error()))
		}
		secretList = append(secretList, secret)
	}
	return secretList, nil
}


// Close закрываем соединение к БД
func (p *psql) Close() error {
	err := p.db.Close()
	if err != nil {
		return fmt.Errorf("can't close DB connection. Error: %s", err.Error())
	}
	return err
}
