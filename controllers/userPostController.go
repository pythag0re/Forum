package controllers

import (
	"forum/db"
	"forum/utils"
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
type CommentWithPost struct {
	PostID     int
	PostTitle  string
	CommentID  int
	Content    string
	CreatedAt  string
	IsOwner    bool
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
func UserCommentsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	loggedIn, sessionUserID := utils.IsAuthenticated(r)

	rows, err := db.DB.Query(`
		SELECT comments.id, comments.content, comments.created_at, posts.id, posts.title
		FROM comments
		JOIN posts ON comments.post_id = posts.id
		WHERE comments.user_id = ?
		ORDER BY comments.created_at DESC
	`, userID)
	if err != nil {
		http.Error(w, "Error loading comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []CommentWithPost
	for rows.Next() {
		var c CommentWithPost
		err := rows.Scan(&c.CommentID, &c.Content, &c.CreatedAt, &c.PostID, &c.PostTitle)
		if err == nil {
			c.IsOwner = loggedIn && sessionUserID == userID
			comments = append(comments, c)
		}
	}

	tmpl := template.Must(template.ParseFiles("_templates_/user_comments.html"))
	tmpl.Execute(w, comments)
}