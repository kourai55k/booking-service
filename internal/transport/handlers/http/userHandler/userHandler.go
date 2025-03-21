package userHandler

import "github.com/kourai55k/booking-service/internal/domain/models"

type UserService interface {
	GetUsers() ([]*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByLogin(login string) (*models.User, error)
	CreateUser(user *models.User) (uint, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type UserHandler struct {
	userService UserService
	logger      Logger
}

func NewUserHandler(userService UserService, logger Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}
