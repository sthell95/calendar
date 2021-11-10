package middleware

import (
	"net/http"

	"calendar.com/pkg/domain/service"
	"calendar.com/pkg/logger"
)

type key string

const AuthService key = "auth"

func Authorization(authService service.Authorization, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		err := authService.IsAuthorized(token)
		if err != nil {
			logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
			return
		}

		next.ServeHTTP(w, r)
	})
}
