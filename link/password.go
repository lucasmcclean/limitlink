package link

import "golang.org/x/crypto/bcrypt"

const hashCost = bcrypt.DefaultCost

// generateHash hashes a plaintext password.
// Returns the resulting hash or an error if hashing fails.
func generateHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(bytes), err
}

// VerifyPassword checks whether the given password matches the hashed password.
// Returns true if the password is correct, false otherwise.
func VerifyPassword(hash string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil, nil
}
