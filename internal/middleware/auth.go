package middleware

import (
	"context"
	"net/http"
	"t-challenge/internal/durable"
	"t-challenge/internal/response"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid, err := durable.ValidateToken(r)
		if err != nil {
			response.WriteResponse(w, &response.Response{
				Status:  http.StatusUnauthorized,
				Success: false,
				Message: "Unauthorized",
			})
			return
		}

		ctx := context.WithValue(r.Context(), "uuid", uuid)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
