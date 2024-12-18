package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
    length := 32
    token := make([]byte, length)
    _, err := rand.Read(token)
    return hex.EncodeToString(token), err
}
