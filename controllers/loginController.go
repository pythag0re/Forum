package controllers

import (
	"database/sql"
	"fmt"
	"forum/db"
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

	var storedPseudo, storedPassword string
	query := "SELECT pseudo, password FROM players WHERE pseudo = ?"
	err = db.DB.QueryRow(query, pseudo).Scan(&storedPseudo, &storedPassword)

	fmt.Println("🔍 cherche le pseudo :", pseudo)

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
		http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
		fmt.Println("Mot de passe incorrect")
		return
	}
	cookie := &http.Cookie{
		Name:     "pseudo",
		Value:    pseudo,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/game-home", http.StatusSeeOther)
	fmt.Println("Connexion réussie, go vers le forum")
}
