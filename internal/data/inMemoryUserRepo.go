package data

import (
	"fmt"
	"sync"

	"github.com/kourai55k/booking-service/internal/domain"
	"github.com/kourai55k/booking-service/internal/domain/models"
)

type InMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[uint]*models.User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[uint]*models.User),
	}
}

func (r *InMemoryUserRepo) GetUsers() ([]*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}

	return users, nil
}

func (r *InMemoryUserRepo) GetUserByID(id uint) (*models.User, error) {
	const op = "InMemoryUserRepo.GetUserById"
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrUserNotFound)
	}

	return user, nil
}

func (r *InMemoryUserRepo) GetUserByLogin(login string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Login == login {
			return u, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (r *InMemoryUserRepo) CreateUser(user *models.User) (uint, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if a user with the same login already exists
	for _, existingUser := range r.users {
		if existingUser.Login == user.Login {
			return 0, fmt.Errorf("InMemoryUserRepo.CreateUser: %w", domain.ErrUserAlreadyExists)
		}
	}

	// Generate a new ID (ensure we don’t overwrite an existing one)
	id := uint(len(r.users) + 1)
	for _, exists := r.users[id]; exists; id++ { // Ensure the ID is unique
	}

	// Set the new user ID and save it
	user.ID = id
	r.users[id] = user

	return id, nil
}

func (r *InMemoryUserRepo) UpdateUser(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return domain.ErrUserNotFound
	}

	r.users[user.ID] = user

	return nil
}

func (r *InMemoryUserRepo) DeleteUser(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.users, id)

	return nil
}
