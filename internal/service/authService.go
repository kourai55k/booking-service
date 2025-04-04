package service

import (
	"errors"
	"fmt"

	"github.com/kourai55k/booking-service/internal/domain"
	"github.com/kourai55k/booking-service/internal/domain/models"
	"github.com/kourai55k/booking-service/pkg/hashing"
	jwthelper "github.com/kourai55k/booking-service/pkg/jwtHelper"
)

type AuthUserRepository interface {
	GetUserByLogin(login string) (*models.User, error)
	CreateUser(user *models.User) (uint, error)
}

type AuthService struct {
	repo AuthUserRepository
}

func NewAuthService(repo AuthUserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(user *models.User) (uint, error) {
	const op = "AuthService.Register"
	id, err := s.repo.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *AuthService) Login(login, password string) (string, error) {
	const op = "AuthService.Login"

	// check if user with provided login exist
	user, err := s.repo.GetUserByLogin(login)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", fmt.Errorf("%s: %w", op, domain.ErrUserNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// check if password is correct
	err = hashing.CheckPassword(user.HashPass, password)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, domain.ErrWrongPassword)
	}

	// generating JWT
	token, err := jwthelper.GenerateToken(user)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
