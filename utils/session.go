package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
