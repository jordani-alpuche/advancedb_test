package validator

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// ValidatePassword checks if the password meets the rules
func ValidatePassword(password string) error {
	
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters long")
    }
    return nil
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashed), nil
}
