package link

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const cost = bcrypt.DefaultCost

// hashPassword hashes a plaintext password.
// Returns the resulting hash or an error if hashing fails.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

// IsCorrectPassword checks whether the given password matches the hashed password
// on the link.
// Returns true if the password is correct or if there is no password and false
// otherwise.
// Returns an error only if there was an unexpected problem verifying the hash.
func (l *Link) IsCorrectPassword(password string) (bool, error) {
	if l.PasswordHash == nil {
		return true, nil
	}
	err := bcrypt.CompareHashAndPassword([]byte(*l.PasswordHash), []byte(password))
	switch err {
	case bcrypt.ErrMismatchedHashAndPassword:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, fmt.Errorf("error comparing password and hash: %w", err)
	}
}
