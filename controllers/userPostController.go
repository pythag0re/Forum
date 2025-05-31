package controllers

import (
	"forum/db"
	"html/template"
	"net/http"
	"strconv"
)

type PublicPost struct {
	ID        int
	Title     string
	Content   string
	CreatedAt string
}

func UserPostsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, title, content, created_at 
		FROM posts 
		WHERE user_id = ? 
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []PublicPost
	for rows.Next() {
		var p PublicPost
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)
		if err == nil {
			posts = append(posts, p)
		}
	}

	tmpl := template.Must(template.ParseFiles("_templates_/user_posts.html"))
	tmpl.Execute(w, posts)
}
