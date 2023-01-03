package handler

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

// mwTokenAuth 	проверяет наличие токена в sessions, если есть прикрепляет userID к
// 				контексту для дальнейшей авторизации действий пользователя
func mwTokenAuth(c *Controller) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// проверяем наличие токена
			var userID int
			token := r.Header.Get("Authorization")
			userID, err = c.sessions.GetUserID(r.Context(), token)
			if err != nil {
				c.log.Info("Can`t get userID by token",
					zap.String("method", "Handler.mwTokenAuth"),
					zap.String("token", token),
					zap.String("error", err.Error()))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			c.log.Debug("Success middleware lookup userID",
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