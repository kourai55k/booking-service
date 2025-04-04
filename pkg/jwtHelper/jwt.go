package jwthelper

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kourai55k/booking-service/internal/domain/models"
)

var (
	secretKey = []byte(os.Getenv("JWT_SECRET"))
)

// Define a struct for custom claims that include Login and Role
type CustomClaims struct {
	Login string `json:"login"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken generates a new JWT token for the user with login and role included
func GenerateToken(user *models.User) (string, error) {
	// Define token expiration time (1 hour in this example)
	expirationTime := time.Now().Add(1 * time.Hour)

	// Create custom claims with login, role, and expiration time
	claims := &CustomClaims{
		Login: user.Login,
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token using HMAC SHA256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseToken parses the JWT token and returns the claims
func ParseToken(tokenString string) (*jwt.Token, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method (ensure it's the correct one)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Return the secret key to verify the token signature
		return secretKey, nil
	})

	return token, err
}
