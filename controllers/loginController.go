package controllers

import (
	"database/sql"
	"fmt"
	"forum/db"
	"forum/utils"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	pseudo := r.FormValue("pseudo")
	password := r.FormValue("password")

	var userID int
	var storedPseudo, storedPassword string

	query := "SELECT id, pseudo, password FROM users WHERE pseudo = ?"
	err = db.DB.QueryRow(query, pseudo).Scan(&userID, &storedPseudo, &storedPassword)

	if err == sql.ErrNoRows {
		fmt.Println("Pseudo introuvable")
		renderLoginWithError(w, "User doesn't exist. Please sign up.")
		return
	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		renderLoginWithError(w, "Incorrect password.")
		return
	}

	sessionToken := utils.GenerateSessionToken()
	_, err = db.DB.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", userID, sessionToken)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// cookie session_token
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
	})

	// cookie user id
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    fmt.Sprintf("%d", userID),
		Path:     "/",
		HttpOnly: false, // Important : JS peut y accéder
	})

	http.Redirect(w, r, "/landing", http.StatusSeeOther)
}

func renderLoginWithError(w http.ResponseWriter, errorMsg string) {
	tmpl := template.Must(template.ParseFiles("_templates_/login.html"))
	data := struct {
		Error string
	}{
		Error: errorMsg,
	}
	tmpl.Execute(w, data)
}
