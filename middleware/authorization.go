package middleware

import (
	"net/http"

	"calendar.com/pkg/response"

	"calendar.com/pkg/domain/service"
	"calendar.com/pkg/logger"
)

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := service.IsAuthorized(r)
		if err != nil {
			logger.NewLogger().Write(logger.Error, err.Error(), "create-event")
			response.NewPrint().PrettyPrint(w, struct {
				Message string
			}{Message: err.Error()}, response.WithCode(http.StatusUnauthorized))
			return
		}

		next.ServeHTTP(w, r)
	})
}
