package repository

import (
	"context"
	"errors"
	"fmt"

	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose"
	"go.uber.org/zap"
)

// Реализация создания нового подклчения к БД

// DB интерфейс для работы с Psql
type DB interface {
	UserExist(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, login string, password string) (int, error)
	///
	SignIn() error
	SignOut() error
	CreateSecret(userID int) error
	UpdateSecret(userID int, id int) error
	ListSecret(userID int) error
	DeleteSecret(userID int, id int) error
	Close() error
}

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
		return 0, fmt.Errorf("create new user fail: %s", err.Error())
	}
	return userID, nil
}

// SignIn логин пользователя
func (p *psql) SignIn() error {
	return nil
}

// SignOut логаут пользователя
func (p *psql) SignOut() error {
	return nil
}

// CreateSecret создать новый секрет
func (p *psql) CreateSecret(userID int) error {
	return nil
}

// UpdateSecret обновить секрет по id
func (p *psql) UpdateSecret(userID int, id int) error {
	return nil
}

// ListSecret список секретов для текущего пользователя
func (p *psql) ListSecret(userID int) error {
	return nil
}

// DeleteSecret удалить секрет
func (p *psql) DeleteSecret(userID int, id int) error {
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
