package middleware

import (
	"context"
	"net/http"

	"calendar.com/pkg/domain/service"
	"calendar.com/pkg/logger"
	"calendar.com/pkg/response"
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

func AuthorizedUserToContext(next http.Handler) http.Handler {
	return http.Handler(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(context.Background())
	})
}
