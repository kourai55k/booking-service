package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/kourai55k/booking-service/internal/domain"
	jwthelper "github.com/kourai55k/booking-service/pkg/jwtHelper"
)

// AdminMiddleware ensures the user is authenticated and is an "admin"
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		// Parse token and get claims
		claims, err := jwthelper.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Validate role
		if claims.Role != "admin" {
			http.Error(w, "forbidden: admin access required", http.StatusForbidden)
			return
		}

		// Add userID and role to context
		ctx := context.WithValue(r.Context(), domain.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, domain.RoleKey, claims.Role)

		// Continue request with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
