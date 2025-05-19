package controllers

import (
	"database/sql"
	"fmt"
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
	hashedPassword, _ := HashPassword(password)

	var existingPseudo string
	queryCheck := "SELECT pseudo FROM users WHERE pseudo = ?"
	err = db.DB.QueryRow(queryCheck, pseudo).Scan(&existingPseudo)

	if err != nil && err != sql.ErrNoRows {
		log.Println("Database error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err == nil {
		fmt.Println("Le pseudo est déjà utilisé. Essayez un autre!")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if password == "" || len(password) < 8 || len(password) > 20 ||
		!strings.ContainsAny(password, "0123456789") ||
		!strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") ||
		!strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		fmt.Println("Mot de passe faible")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	queryInsert := `INSERT INTO users (email, pseudo, password) VALUES (?, ?, ?)`
	_, err = db.DB.Exec(queryInsert, email, pseudo, hashedPassword)
	if err != nil {
		log.Println("Erreur lors de l'enregistrement :", err)
		http.Error(w, "Erreur lors de la création du compte", http.StatusInternalServerError)
		return
	}

	fmt.Println("Inscription réussie !")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
