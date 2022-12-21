package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)


func NewRouter(c *Controller, log *zap.Logger) http.Handler {
	r := chi.NewRouter()

	// Аутентификация
	// Роутинг
	r.Use(mwTokenAuth())
	r.Get("/version", c.Version)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/version", c.Version)
	})
	return r
}