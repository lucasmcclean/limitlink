package link

import (
	"crypto/rand"
	"math/big"
)

const (
	minSlugLen = 6
	maxSlugLen = 12

	lettersCharset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numericCharset      = "0123456789"
	alphanumericCharset = lettersCharset + numericCharset
)

var reservedPrefixes = []string{
	"static", "admin", "links", "stat",
}

// generateSlug creates a cryptographically secure random string (slug) of the specified length,
// using the specified charset: "letters", "numbers", or "alphanumeric".
// If charset is empty or unrecognized, "alphanumeric" is used by default.
func generateSlug(length int, charset string) (string, error) {
	var chars string
	switch charset {
	case "letters":
		chars = lettersCharset
	case "numbers":
		chars = numericCharset
	case "alphanumeric":
		chars = alphanumericCharset
	default:
		chars = alphanumericCharset
	}

	slug := make([]byte, length)
	charsetLen := big.NewInt(int64(len(chars)))
	for {
		for i := range length {
			n, err := rand.Int(rand.Reader, charsetLen)
			if err != nil {
				return "", err
			}
			slug[i] = chars[n.Int64()]
		}
		if !isReserved(string(slug)) {
			break
		}
	}
	return string(slug), nil
}

func isReserved(slug string) bool {
	for _, prefix := range reservedPrefixes {
		if len(slug) >= len(prefix) && slug[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}
