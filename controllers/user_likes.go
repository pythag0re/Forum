package controllers

import (
	"forum/db"
	"html/template"
	"net/http"
	"strconv"
)

type LikedPost struct {
	ID        int
	Title     string
	Content   string
	Author    string
	CreatedAt string
}

func UserLikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT posts.id, posts.title, posts.content, users.pseudo, posts.created_at
		FROM likes
		JOIN posts ON likes.post_id = posts.id
		JOIN users ON posts.user_id = users.id
		WHERE likes.user_id = ?
		ORDER BY posts.created_at DESC
	`, userID)

	if err != nil {
		http.Error(w, "Error fetching liked posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var likedPosts []LikedPost
	for rows.Next() {
		var post LikedPost
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt)
		if err == nil {
			likedPosts = append(likedPosts, post)
		}
	}

	tmpl := template.Must(template.ParseFiles("_templates_/user_likes.html"))
	tmpl.Execute(w, likedPosts)
}
