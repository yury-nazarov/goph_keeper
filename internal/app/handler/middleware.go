package handler

import "net/http"

// mwTokenAuth проверяет токен в
func mwTokenAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// TODO: проверяем наличие токена
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}