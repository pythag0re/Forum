package controllers

import (
	"forum/db"
	"forum/utils"
	"html/template"
	"net/http"
)

type Comment struct {
  Author    string
  Text      string
  CreatedAt string
}

type PostDetail struct {
  ID        int
  Title     string
  Content   string
  Author    string
  CreatedAt string
  Comments  []Comment
}

func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
  id := r.URL.Query().Get("id")
  if id == "" {
    http.Error(w, "Missing post ID", http.StatusBadRequest)
    return
  }

  var post PostDetail
  err := db.DB.QueryRow(`
    SELECT posts.id, posts.title, posts.content, users.pseudo, posts.created_at
    FROM posts
    JOIN users ON posts.user_id = users.id
    WHERE posts.id = ?
  `, id).Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt)
  if err != nil {
    http.Error(w, "Post not found", http.StatusNotFound)
    return
  }

  rows, err := db.DB.Query(`
    SELECT users.pseudo, comments.content, comments.created_at
    FROM comments
    JOIN users ON comments.user_id = users.id
    WHERE comments.post_id = ?
    ORDER BY comments.created_at ASC
  `, id)

  if err == nil {
    defer rows.Close()
    for rows.Next() {
      var c Comment
      rows.Scan(&c.Author, &c.Text, &c.CreatedAt)
      post.Comments = append(post.Comments, c)
    }
  }

  tmpl := template.Must(template.ParseFiles("_templates_/post_detail.html"))
  tmpl.Execute(w, post)
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