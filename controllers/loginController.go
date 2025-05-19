package controllers

import (
	"database/sql"
	"fmt"
	"forum/db"
	"forum/utils"
	"log"
	"net/http"
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
	fmt.Println("Recherche du pseudo :", pseudo)

	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		fmt.Println("Utilisateur introuvable. Redirection vers l'inscription.")
		return
	} else if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if storedPassword != password {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		fmt.Println("Mot de passe incorrect")
		return
	}

	sessionToken := utils.GenerateSessionToken()

	
	_, err = db.DB.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", userID, sessionToken)
	if err != nil {
		log.Println("Erreur lors de l'insertion du token :", err)
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
	fmt.Println("Connexion réussie. Redirection vers le forum.")
}

