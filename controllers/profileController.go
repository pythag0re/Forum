package controllers

import (
	"forum/db"
	"forum/utils"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)


func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := utils.GetUserIDFromSession(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	oldPassword := r.FormValue("oldPassword")
	newPassword := r.FormValue("newPassword")
	confirmPassword := r.FormValue("confirmPassword")

	if newPassword != confirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	var hashedPassword string
	err = db.DB.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(oldPassword))
	if err != nil {
		http.Error(w, "Incorrect current password", http.StatusUnauthorized)
		return
	}

	newHashed, _ := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	_, err = db.DB.Exec("UPDATE users SET password = ? WHERE id = ?", string(newHashed), userID)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	_, err = db.DB.Exec("DELETE FROM sessions WHERE token = ?", cookie.Value)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	expiredCookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w, expiredCookie)

	http.Redirect(w, r, "/landing", http.StatusSeeOther)
}

func DeleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ok, userID := utils.IsAuthenticated(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	_, err := db.DB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		log.Println("Erreur lors de la suppression du profil:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		log.Println("Erreur lors de la suppression de la session:", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/register", http.StatusSeeOther)
}

func MyPostsHandler(w http.ResponseWriter, r *http.Request) {
	ok, userID := utils.IsAuthenticated(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, title, content, created_at 
		FROM posts 
		WHERE user_id = ? 
		ORDER BY created_at DESC`, userID)
	if err != nil {
		log.Println("Erreur récupération de mes posts :", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []PostData
	for rows.Next() {
		var p PostData
		p.Author = "You" // facultatif
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)
		if err == nil {
			posts = append(posts, p)
		}
	}

	tmpl := template.Must(template.ParseFiles("_templates_/my_posts.html"))
	tmpl.Execute(w, posts)
}
