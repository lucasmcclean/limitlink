package link

import (
	"crypto/rand"
	"math/big"
)

const adminTokenLen = 22

// generateAdminToken returns a securely generated random alphanumeric string.
// The token is used as an admin secret and uses crypto/rand for secure randomness.
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
