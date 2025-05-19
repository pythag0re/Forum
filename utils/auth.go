package utils

import (
	"forum/db"
	"net/http"
)

func IsAuthenticated(r *http.Request) (bool, int) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false, 0
	}

	var userID int
	err = db.DB.QueryRow("SELECT user_id FROM sessions WHERE token = ?", cookie.Value).Scan(&userID)
	if err != nil {
		return false, 0
	}

	return true, userID
}
