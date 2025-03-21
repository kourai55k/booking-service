package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kourai55k/booking-service/internal/domain/models"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// CreateUserTable creates the "users" table if it doesn't exist.
func (u *UserRepo) CreateUserTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		login TEXT UNIQUE NOT NULL,
		hashpass TEXT NOT NULL,
		role TEXT
	);
	`
	_, err := u.pool.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("CreateUserTable: failed to create table: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by its ID.
func (r *UserRepo) GetUserByID(id uint) (*models.User, error) {
	query := "SELECT id, name, login, hashpass FROM users WHERE id = $1"
	row := r.pool.QueryRow(context.Background(), query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Login, &user.HashPass); err != nil {
		// TODO: check if user not found and return domain.ErrUserNotFound
		return nil, fmt.Errorf("UserRepo.GetUserByID: %w", err)
	}

	return &user, nil
}

// CreateUser creates a new user in the database and returns the new user's id.
func (r *UserRepo) CreateUser(user *models.User) (uint, error) {
	query := "INSERT INTO users (name, login, hashpass) VALUES ($1, $2, $3) RETURNING id"
	var id uint
	err := r.pool.QueryRow(context.Background(), query, user.Name, user.Login, user.HashPass).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("UserRepo.CreateUser: %w", err)
	}
	return id, nil
}

// GetUsers retrieves all users from the database.
func (r *UserRepo) GetUsers() ([]*models.User, error) {
	query := "SELECT id, name, login FROM users"
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetUsers: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Login); err != nil {
			return nil, fmt.Errorf("UserRepo.GetUsers: %w", err)
		}
		users = append(users, &user)
	}

	return users, nil
}

// GetUserByLogin retrieves a user by its login.
func (r *UserRepo) GetUserByLogin(login string) (*models.User, error) {
	query := "SELECT id, name, login FROM users WHERE login = $1"
	row := r.pool.QueryRow(context.Background(), query, login)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Login); err != nil {
		return nil, fmt.Errorf("UserRepo.GetUserByLogin: %w", err)
	}

	return &user, nil
}

// UpdateUser updates an existing user in the database.
func (r *UserRepo) UpdateUser(user *models.User) error {
	query := "UPDATE users SET name = $2, login = $3 WHERE id = $1"
	_, err := r.pool.Exec(context.Background(), query, user.ID, user.Name, user.Login)
	return err
}

// DeleteUser deletes a user in the database.
func (r *UserRepo) DeleteUser(id uint) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.pool.Exec(context.Background(), query, id)
	return err
}
