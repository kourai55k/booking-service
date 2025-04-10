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
	tokenTTL  = 1 * time.Hour
)

// CustomClaims includes additional fields for user identity
type CustomClaims struct {
	UserID uint   `json:"user_id"`
	Login  string `json:"login"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken creates a signed JWT with custom claims
func GenerateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(tokenTTL)

	claims := &CustomClaims{
		UserID: user.ID,
		Login:  user.Login,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseToken parses and validates JWT, returning custom claims
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
