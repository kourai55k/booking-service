package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUsersNotFound     = errors.New("users not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrWrongPassword     = errors.New("wrong password")
)
