package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/kourai55k/booking-service/internal/domain"
	jwthelper "github.com/kourai55k/booking-service/pkg/jwtHelper"
)

// AuthMiddleware is a middleware function that checks if the request has a valid JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Split the header into "Bearer token"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
			return
		}

		// Get the token from the second part
		tokenStr := parts[1]

		// Parse token and get claims
		claims, err := jwthelper.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Add userID and role to context
		ctx := context.WithValue(r.Context(), domain.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, domain.RoleKey, claims.Role)

		// Continue request with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
