package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kourai55k/booking-service/internal/domain"
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

// CreateUser creates a new user in the database and returns the new user's id.
func (r *UserRepo) CreateUser(user *models.User) (uint, error) {
	query := "INSERT INTO users (name, login, hashpass, role) VALUES ($1, $2, $3, $4) RETURNING id"
	var id uint
	err := r.pool.QueryRow(context.Background(), query, user.Name, user.Login, user.HashPass, user.Role).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Check unique constraint violation
			return 0, fmt.Errorf("UserRepo.CreateUser: %w", domain.ErrUserAlreadyExists)
		}
		return 0, fmt.Errorf("UserRepo.CreateUser: %w", err)
	}
	return id, nil
}

// GetUsers retrieves all users from the database.
func (r *UserRepo) GetUsers() ([]*models.User, error) {
	query := "SELECT id, name, login, hashpass, role FROM users"
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetUsers: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.HashPass, &user.Role); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("UserRepo.GetUsers: %w", domain.ErrUsersNotFound)
			}
			return nil, fmt.Errorf("UserRepo.GetUsers: %w", err)
		}
		users = append(users, &user)
	}

	return users, nil
}

// GetUserByID retrieves a user by its ID.
func (r *UserRepo) GetUserByID(id uint) (*models.User, error) {
	query := "SELECT id, name, login, hashpass, role FROM users WHERE id = $1"
	row := r.pool.QueryRow(context.Background(), query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Login, &user.HashPass, &user.Role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("UserRepo.GetUserByID: %w", domain.ErrUserNotFound)
		}
		return nil, fmt.Errorf("UserRepo.GetUserByID: %w", err)
	}

	return &user, nil
}

// GetUserByLogin retrieves a user by its login.
func (r *UserRepo) GetUserByLogin(login string) (*models.User, error) {
	query := "SELECT id, name, login, hashpass, role FROM users WHERE login = $1"
	row := r.pool.QueryRow(context.Background(), query, login)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Login, &user.HashPass, &user.Role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("UserRepo.GetUserByLogin: %w", domain.ErrUserNotFound)
		}
		return nil, fmt.Errorf("UserRepo.GetUserByLogin: %w", err)
	}

	return &user, nil
}

// UpdateUser updates an existing user in the database.
func (r *UserRepo) UpdateUser(user *models.User) error {
	// Проверяем, существует ли пользователь
	existsQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)"
	var exists bool
	err := r.pool.QueryRow(context.Background(), existsQuery, user.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("UserRepo.UpdateUser: check user existence failed: %w", err) // Ошибка БД
	}

	// Если пользователя нет, возвращаем ошибку
	if !exists {
		return domain.ErrUserNotFound
	}

	// Выполняем UPDATE только если пользователь существует
	query := `
		UPDATE users
		SET name = COALESCE(NULLIF($2, ''), name),
		    name = COALESCE(NULLIF($3, ''), name),
		    login = COALESCE(NULLIF($4, ''), login),
		    hashpass = COALESCE(NULLIF($5, ''), hashpass),
		    role = COALESCE(NULLIF($6, ''), role)
		WHERE id = $1
	`
	_, err = r.pool.Exec(context.Background(), query, user.ID, user.Name, user.Login, user.HashPass, user.Role)
	if err != nil {
		return fmt.Errorf("UpdateUser failed: %w", err)
	}

	return nil
}

// DeleteUser deletes a user from the database.
func (r *UserRepo) DeleteUser(id uint) error {
	// Проверяем, существует ли пользователь
	existsQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)"
	var exists bool
	err := r.pool.QueryRow(context.Background(), existsQuery, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("UserRepo.DeleteUser: check user existence failed: %w", err) // Ошибка БД
	}

	// Если пользователя нет, возвращаем ошибку
	if !exists {
		return domain.ErrUserNotFound
	}

	// Удаляем пользователя
	query := "DELETE FROM users WHERE id = $1"
	_, err = r.pool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("UserRepo.DeleteUser failed: %w", err)
	}

	return nil
}
