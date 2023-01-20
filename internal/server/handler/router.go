package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(authController *authController, secretController *secretController, msController *msController) http.Handler {
	r := chi.NewRouter()

	// Роутинг
	r.Get("/version", msController.Version)
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", authController.SignUp)
			r.Post("/signin", authController.SignIn)
			r.Group(func(r chi.Router) {
				r.Use(mwTokenAuth(authController))
				r.Delete("/signout", authController.SignOut)
			})
		})
		r.Route("/secret", func(r chi.Router) {
			r.Use(mwTokenAuth(authController))
			r.Post("/new", secretController.SecretNew)
			r.Get("/list", secretController.SecretList)
			r.Get("/{secretID}", secretController.GetSecretByID)
			r.Put("/update", secretController.UpdateSecretByID)
			r.Delete("/delete/{secretID}", secretController.DeleteSecretByID)
		})
	})
	return r
}
