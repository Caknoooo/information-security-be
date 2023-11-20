package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateKey(len int) (string, error) {
	bytes := make([]byte, len)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	newKey := hex.EncodeToString(bytes)
	return newKey, nil
}
