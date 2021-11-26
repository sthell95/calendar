package middleware

import (
	"context"
	"net/http"

	"calendar.com/pkg/domain/service"
	"calendar.com/pkg/logger"
	"calendar.com/pkg/response"
)

type CurrentUser string

const UserId CurrentUser = "currentUserId"

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := service.Validate(r)
		if err != nil {
			logger.NewLogger().Write(logger.Error, err.Error(), "authorization")
			response.NewPrint().PrettyPrint(w, struct {
				Message string
			}{Message: err.Error()}, response.WithCode(http.StatusUnauthorized))

			return
		}
		ctx := context.WithValue(r.Context(), UserId, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
