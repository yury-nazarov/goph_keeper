package repository

import "go.uber.org/zap"

// Реализация Inmemory Cache в виде мапы

// Sessions - интерфейс для работы с кешем
type Sessions interface {
	AddToken() error
	GetToken() error
	Close()
}

type inmemorySessionStorage struct {
	store map[int]string
	log   zap.Logger
}

func NewSessions(log *zap.Logger) (*inmemorySessionStorage, error) {
	ic := &inmemorySessionStorage{
		store: make(map[int]string),
		log:   *log,
	}
	log.Info("Sessions storage init success")
	return ic, nil
}

func (c *inmemorySessionStorage) AddToken() error {
	return nil
}

func (c *inmemorySessionStorage) GetToken() error {
	return nil
}

func (c *inmemorySessionStorage) Close() {
	c.log.Info("close inmemory cache storage")
	c.store = nil
}
