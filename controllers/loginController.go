package controllers

import (
	"database/sql"
	"forum/db"
	"forum/utils"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Données de formulaire invalides", http.StatusBadRequest)
		return
	}

	pseudo := r.FormValue("pseudo")
	password := r.FormValue("password")

	var userID int
	var storedPseudo, storedPassword string

	query := "SELECT id, pseudo, password FROM users WHERE pseudo = ?"
	err := db.DB.QueryRow(query, pseudo).Scan(&userID, &storedPseudo, &storedPassword)

	if err == sql.ErrNoRows {
		renderLoginWithError(w, "Utilisateur introuvable. Veuillez vous inscrire.")
		return
	} else if err != nil {
		http.Error(w, "Erreur serveur interne", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		renderLoginWithError(w, "Mot de passe incorrect.")
		return
	}

	sessionToken := utils.GenerateSessionToken()
	_, err = db.DB.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", userID, sessionToken)
	if err != nil {
		http.Error(w, "Erreur de session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
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
