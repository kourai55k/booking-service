package hashing

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(hashPass, password string) error {
	// Compare the hashed password with the plain password
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
	if err != nil {
		// If there was an error comparing (usually due to mismatched passwords), return false
		return err
	}
	return nil
}
