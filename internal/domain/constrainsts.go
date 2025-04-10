package domain

// contextKey is a custom type for context keys
type contextKey string

// context keys
const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

// validation rules
const (
	MinPasswordLength = 8
)
