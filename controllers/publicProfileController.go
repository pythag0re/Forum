package controllers

import (
	"forum/db"
	"html/template"
	"net/http"
	"forum/utils"
	"strconv"
	"strings"
)


func PublicProfileHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/user/")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	authenticated, currentUserID := utils.IsAuthenticated(r)
	isOwner := authenticated && currentUserID == userID

	var pseudo, photo, createdAt string
	var postCount, commentCount int

	err = db.DB.QueryRow(`
		SELECT pseudo, profile_picture, created_at 
		FROM users 
		WHERE id = ?`, userID).Scan(&pseudo, &photo, &createdAt)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusInternalServerError)
		return
	}

	db.DB.QueryRow(`SELECT COUNT(*) FROM posts WHERE user_id = ?`, userID).Scan(&postCount)
	db.DB.QueryRow(`SELECT COUNT(*) FROM comments WHERE user_id = ?`, userID).Scan(&commentCount)

	if photo == "" {
		photo = "avatar.png"
	}

	data := struct {
		ID            int
		Pseudo        string
		ProfilePic    string
		MemberSince   string
		PostCount     int
		CommentCount  int
		IsOwner       bool
	}{
		ID:            userID,
		Pseudo:        pseudo,
		ProfilePic:    photo,
		MemberSince:   createdAt,
		PostCount:     postCount,
		CommentCount:  commentCount,
		IsOwner:       isOwner,
	}

	tmpl := template.Must(template.ParseFiles("_templates_/user_profile.html"))
	tmpl.Execute(w, data)
}