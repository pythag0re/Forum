package controllers

import (
	"database/sql"
	"forum/db"
	"html/template"
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
	confirmPassword := r.FormValue("confirm_password")

	if password != confirmPassword {
		renderRegisterWithError(w, "Passwords do not match.")
		return
	}

	if len(password) < 8 || len(password) > 20 ||
		!strings.ContainsAny(password, "0123456789") ||
		!strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") ||
		!strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		renderRegisterWithError(w, "Password must be 8–20 characters long and include upper, lower and digit.")
		return
	}

	var existing string
	err = db.DB.QueryRow(`SELECT pseudo FROM users WHERE pseudo = ?`, pseudo).Scan(&existing)
	if err == nil {
		renderRegisterWithError(w, "Username already taken.")
		return
	} else if err != sql.ErrNoRows {
		log.Println("DB error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	hashed, err := HashPassword(password)
	if err != nil {
		renderRegisterWithError(w, "Internal error during password hashing.")
		return
	}

	_, err = db.DB.Exec(`INSERT INTO users (email, pseudo, password) VALUES (?, ?, ?)`, email, pseudo, hashed)
	if err != nil {
		log.Println("DB insert error:", err)
		renderRegisterWithError(w, "Could not create account.")
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func renderRegisterWithError(w http.ResponseWriter, errorMsg string) {
	tmpl := template.Must(template.ParseFiles("_templates_/register.html"))
	data := struct {
		Error string
	}{
		Error: errorMsg,
	}
	tmpl.Execute(w, data)
}
