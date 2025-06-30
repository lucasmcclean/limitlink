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

// VerifyPassword checks whether the given password matches the hashed password.
// Returns true if the password is correct, false otherwise.
// Returns an error only if there was an unexpected problem verifying the hash.
func VerifyPassword(hash string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	switch err {
	case bcrypt.ErrMismatchedHashAndPassword:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, fmt.Errorf("error comparing password and hash: %w", err)
	}
}
