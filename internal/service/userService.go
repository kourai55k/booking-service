package service

import (
	"fmt"

	"github.com/kourai55k/booking-service/internal/domain/models"
)

type UserRepository interface {
	GetUsers() ([]*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByLogin(login string) (*models.User, error)
	CreateUser(user *models.User) (uint, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers() ([]*models.User, error) {
	const op = "UserService.GetUsers"

	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	const op = "UserService.GetUserById"

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UserService) GetUserByLogin(login string) (*models.User, error) {
	const op = "UserService.GetUserByLogin"

	user, err := s.repo.GetUserByLogin(login)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UserService) CreateUser(user *models.User) (uint, error) {
	const op = "UserService.CreateUser"

	id, err := s.repo.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	const op = "UserService.UpdateUser"

	err := s.repo.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *UserService) DeleteUser(id uint) error {
	const op = "UserService.DeleteUser"

	err := s.repo.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
