package controllers

import (
	"forum/db"
	"forum/utils"
	"html/template"
	"net/http"
)


func MyCommentsHandler(w http.ResponseWriter, r *http.Request) {
	ok, userID := utils.IsAuthenticated(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := db.DB.Query(`
		SELECT posts.id, posts.title, comments.id, comments.content, comments.created_at
		FROM comments
		JOIN posts ON comments.post_id = posts.id
		WHERE comments.user_id = ?
		ORDER BY comments.created_at DESC`, userID)
	if err != nil {
		http.Error(w, "Error loading your comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type CommentWithPost struct {
		PostID     int
		PostTitle  string
		CommentID  int
		Content    string
		CreatedAt  string
	}

	var comments []CommentWithPost
	for rows.Next() {
		var c CommentWithPost
		err := rows.Scan(&c.PostID, &c.PostTitle, &c.CommentID, &c.Content, &c.CreatedAt)
		if err == nil {
			comments = append(comments, c)
		}
	}

	tmpl := template.Must(template.ParseFiles("_templates_/my_comments.html"))
	tmpl.Execute(w, comments)
}


func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/landing", http.StatusSeeOther)
		return
	}

	ok, userID := utils.IsAuthenticated(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	postID := r.FormValue("post_id")
	content := r.FormValue("comment")

	if postID == "" || content == "" {
		http.Error(w, "Missing comment content or post ID", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`
		INSERT INTO comments (user_id, post_id, content, created_at) 
		VALUES (?, ?, ?, datetime('now'))`, userID, postID, content)

	if err != nil {
		http.Error(w, "Unable to save comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
}

func EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		var content string
		err := db.DB.QueryRow("SELECT content FROM comments WHERE id = ?", id).Scan(&content)
		if err != nil {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		data := struct {
			ID      string
			Content string
		}{ID: id, Content: content}

		tmpl := template.Must(template.ParseFiles("_templates_/edit_comment.html"))
		tmpl.Execute(w, data)
	} else if r.Method == http.MethodPost {
		http.Redirect(w, r, "/my-comments", http.StatusSeeOther)
	}
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		_, err := db.DB.Exec("DELETE FROM comments WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/my-comments", http.StatusSeeOther)
	}
}

func UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		content := r.FormValue("content")
		_, err := db.DB.Exec("UPDATE comments SET content = ? WHERE id = ?", content, id)
		if err != nil {
			http.Error(w, "Failed to update comment", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/my-comments", http.StatusSeeOther)
	}
}
