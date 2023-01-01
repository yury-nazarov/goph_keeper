package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

)

func NewRouter(c *Controller) http.Handler {
	r := chi.NewRouter()

	// Роутинг
	r.Get("/version", c.Version)
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", c.SignUp)
			r.Post("/signin", c.SignIn)
			r.Group(func(r chi.Router) {
				r.Use(mwTokenAuth(c))
				r.Delete("/signout", c.SignOut)
			})
		})
		r.Route("/secrets", func(r chi.Router) {
			// TODO
		})



	})
	return r
}
