package repository

import "go.uber.org/zap"

// Реализация Inmemory Cache в виде мапы

// Cache - интерфейс для работы с кешем
type Cache interface {
	AddToken() error
	GetToken() error
	Close()
}

type inmemoryCache struct {
	store map[int]string
	log   zap.Logger
}

func NewCache(log *zap.Logger) (*inmemoryCache, error) {
	ic := &inmemoryCache{
		store: make(map[int]string),
		log:   *log,
	}
	return ic, nil
}

func (c *inmemoryCache) AddToken() error {
	return nil
}

func (c *inmemoryCache) GetToken() error {
	return nil
}

func (c *inmemoryCache) Close() {
	c.log.Info("close inmemory cache storage")
	c.store = nil
}
