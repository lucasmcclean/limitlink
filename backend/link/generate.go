package link

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

var ErrUnrecognizedCharset = errors.New("unrecognized character set")

const (
	alphanumeric = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers      = "0123456789"
)

// generateAdminToken returns a securely generated random alphanumeric string.
// The token is used as an admin secret and uses crypto/rand for secure randomness.
func generateAdminToken(length int) (string, error) {
	token := make([]byte, length)

	charsLen := big.NewInt(int64(len(alphanumeric)))
	for i := range length {
		n, err := rand.Int(rand.Reader, charsLen)
		if err != nil {
			return "", fmt.Errorf("error generating the admin token: %w", err)
		}
		token[i] = alphanumeric[n.Int64()]
	}

	return string(token), nil
}

// generateSlug creates a cryptographically secure random string (slug) of the specified length,
// using the specified charset: "letters", "numbers", or "alphanumeric".
func generateSlug(length int, charset string) (string, error) {
	var chars string
	switch charset {
	case "letters":
		chars = letters
	case "numbers":
		chars = numbers
	case "alphanumeric":
		chars = alphanumeric
	default:
		return "", ErrUnrecognizedCharset
	}

	slug := make([]byte, length)
	charsLen := big.NewInt(int64(len(chars)))
	for i := range length {
		n, err := rand.Int(rand.Reader, charsLen)
		if err != nil {
			return "", fmt.Errorf("error generating the slug: %w", err)
		}
		slug[i] = chars[n.Int64()]
	}
	return string(slug), nil
}
