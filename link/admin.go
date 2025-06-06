package link

import (
	"crypto/rand"
	"math/big"
)

const adminTokenLen  = 22

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
