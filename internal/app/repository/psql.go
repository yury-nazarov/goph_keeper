package repository

import "go.uber.org/zap"

// Реализация создания нового подклчения в PSQL

// DB интерфейс для работы с Psql
type DB interface {
	SignUp() error
	SignIn() error
	SignOut() error
	CreateSecret() error
	UpdateSecret() error
	ListSecret() error
	DeleteSecret() error
	Close()
}

type psql struct {
	log *zap.Logger
}

func NewPostgres(log *zap.Logger) (*psql, error) {
	p := &psql{
		log: log,
	}
	log.Info("DB storage init success")
	return p, nil
}

func (db *psql) SignUp() error {
	return nil
}

func (db *psql) SignIn() error {
	return nil
}

func (db *psql) SignOut() error {
	return nil
}

func (db *psql) CreateSecret() error {
	return nil
}

func (db *psql) UpdateSecret() error {
	return nil
}

func (db *psql) ListSecret() error {
	return nil
}

func (db *psql) DeleteSecret() error {
	return nil
}

func (db *psql) Close() {
}
