package repository

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"sync"
)

// Реализация Inmemory Cache в виде мапы.

// Sessions - интерфейс для работы с кешем
type Sessions interface {
	AddToken(ctx context.Context, token string, userID int) error
	GetUserID(ctx context.Context, token string) (int, error)
	DeleteToken(ctx context.Context, token string) error
	Close() error
}

type inmemorySessionStorage struct {
	store map[string]int
	mu    sync.RWMutex
	log   zap.Logger
}

// NewSessions создает inmemory cache для хранения токенов залогиненых пользователей
func NewSessions(log *zap.Logger) (*inmemorySessionStorage, error) {
	ic := &inmemorySessionStorage{
		store: make(map[string]int),
		log:   *log,
	}
	log.Info("Sessions storage init success")
	return ic, nil
}

// AddToken добавить token и соответствующий ему userID в кеш сессий
func (c *inmemorySessionStorage) AddToken(ctx context.Context, token string, userID int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[token] = userID
	return nil
}

// GetUserID получить userID по token
func (c *inmemorySessionStorage) GetUserID(ctx context.Context, token string) (int, error) {
	userID, ok := c.store[token]
	if !ok {
		return 0, fmt.Errorf("session for token: %s not found", token)
	}
	return userID, nil
}

// DeleteToken удаляет токен из активных сессий
func (c *inmemorySessionStorage) DeleteToken(ctx context.Context, token string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, token)
	return nil
}

// Close завершает работу с inmemory cache
func (c *inmemorySessionStorage) Close() error {
	c.log.Info("close inmemory cache storage")
	c.store = nil
	return nil
}
