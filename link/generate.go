package link

import (
	"crypto/rand"
	"math/big"
)

type slugCharset int

const (
	alphanumeric slugCharset = iota
	letters
	numbers
)

const (
	lettersCharset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numericCharset      = "0123456789"
	alphanumericCharset = lettersCharset + numericCharset
)

const (
	minSlugLen     = 6
	maxSlugLen     = 12
	defaultSlugLen = 7
	adminTokenLen  = 22
)

func generateSlug(length int, charset slugCharset) (string, error) {
	var chars string
	switch charset {
	case letters:
		chars = lettersCharset
	case numbers:
		chars = numericCharset
	case alphanumeric:
		chars = alphanumericCharset
	default:
		chars = alphanumericCharset
	}

	slug := make([]byte, length)
	charsetLen := big.NewInt(int64(len(chars)))
	for i := range length {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		slug[i] = chars[n.Int64()]
	}

	return string(slug), nil
}

func generateAdminToken() (string, error) {
	token := make([]byte, adminTokenLen)

	charsetLen := big.NewInt(int64(len(alphanumericCharset)))
	for i := range adminTokenLen {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		token[i] = alphanumericCharset[n.Int64()]
	}

	return string(token), nil
}
