package controllers

import (
	"fmt"
	"forum/db"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profil", http.StatusSeeOther)
		return
	}

	userID := getAuthenticatedUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username := html.EscapeString(r.FormValue("username"))
	email := html.EscapeString(r.FormValue("email"))
	file, header, err := r.FormFile("profilePic")

	var imagePath string

	if err == nil && header.Filename != "" {
		defer file.Close()
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("user_%d%s", userID, ext)
		imagePath = filepath.Join("_templates_/public/imgs/", filename)

		out, err := os.Create(imagePath)
		if err != nil {
			http.Error(w, "Erreur lors de l'enregistrement de l'image", http.StatusInternalServerError)
			return
		}
		defer out.Close()
		io.Copy(out, file)
	}

	err = db.UpdateUserProfile(userID, email, username)
	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour", http.StatusInternalServerError)
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("Profil mis à jour :", username)
	http.Redirect(w, r, "/profil", http.StatusSeeOther)
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	userID := getAuthenticatedUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := db.DeleteUserByID(userID)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression du profil", http.StatusInternalServerError)
		fmt.Println("Erreur suppression :", err)
		return
	}

	fmt.Println("Utilisateur supprimé avec succès :", userID)
	http.Redirect(w, r, "/landing", http.StatusSeeOther)
}

func getAuthenticatedUserID(r *http.Request) int {
	return 1
}
