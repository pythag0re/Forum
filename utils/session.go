package utils

import (
	"crypto/rand"
	"encoding/hex"
	"forum/db"
)

func GenerateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GetUserIDFromSession(token string) (int, error) {
	var userID int
	err := db.DB.QueryRow("SELECT user_id FROM sessions WHERE token = ?", token).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

