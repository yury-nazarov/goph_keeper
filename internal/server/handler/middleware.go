package handler

import (
	"context"

	"net/http"

	"go.uber.org/zap"
)

// mwTokenAuth 	проверяет наличие токена в sessions, если есть прикрепляет userID к
// 				контексту для дальнейшей авторизации действий пользователя
func mwTokenAuth(authController *authController) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// проверяем наличие токена
			var userID int
			token := r.Header.Get("Authorization")
			userID, err = authController.sessions.GetUserID(r.Context(), token)
			if err != nil {
				authController.log.Info("Can`t get userID by token",
					zap.String("method", "Handler.mwTokenAuth"),
					zap.String("token", token),
					zap.String("error", err.Error()))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			authController.log.Debug("Success middleware lookup userID",
				zap.String("method", "Handler.mwTokenAuth"),
				zap.String("token", token),
				zap.Int("userID", userID))

			// Добавляем userID к контексту
			// далее в методах можно получить доступ через: ctx.Value("userID")
			req := r.WithContext(context.WithValue(r.Context(), "userID", userID))
			*r = *req

			// Передаем request далее
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
