package controllers

import (
	"forum/db"
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

