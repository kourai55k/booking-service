package authHandler

import "github.com/kourai55k/booking-service/internal/domain/models"

type AuthService interface {
	Register(user *models.User) (uint, error)
	Login(login, password string) (string, error)
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type AuthHandler struct {
	authService AuthService
	logger      Logger
}

func NewAuthHandler(authService AuthService, logger Logger) *AuthHandler {
	return &AuthHandler{authService: authService, logger: logger}
}
