package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
		tokenString := parts[1]

		// Parse and validate the token using the helper function
		token, err := jwthelper.ParseToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %s", err.Error()), http.StatusUnauthorized)
			return
		}

		// Check if the token is valid and not expired
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Attach the claims (user information) to the request context (via headers or context)
		// You can use the context to store these values if you prefer not to use headers
		r.Header.Set("UserLogin", claims["login"].(string))
		r.Header.Set("UserRole", claims["role"].(string))

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
