package controllers

import (
	"forum/db"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	pseudo := r.FormValue("pseudo")
	password := r.FormValue("password")
	if len(password) < 8 || len(password) > 20 ||
		!strings.ContainsAny(password, "0123456789") ||
		!strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") ||
		!strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		http.Error(w, "Weak password. Must be 8–20 chars, include upper, lower, and number.", http.StatusBadRequest)
		return
	}

	
	var existing string
	err = db.DB.QueryRow(`SELECT pseudo FROM users WHERE pseudo = ?`, pseudo).Scan(&existing)
	if err == nil {
		http.Error(w, "Username already taken.", http.StatusConflict)
		return
	}

	hashed, err := HashPassword(password)
	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		log.Println("Hash error:", err)
		return
	}

	_, err = db.DB.Exec(`
		INSERT INTO users (email, pseudo, password, profile_picture) 
		VALUES (?, ?, ?, ?)`, email, pseudo, hashed, "")
	if err != nil {
		log.Println("DB insert error:", err)
		http.Error(w, "Error creating account.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
