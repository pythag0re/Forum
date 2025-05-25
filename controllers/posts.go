package controllers

import (
	"forum/db"
	"html/template"
	"log"
	"net/http"
)

type PostData struct {
	ID        int
	Title     string
	Content   string 
	Author    string
	CreatedAt string
	Likes     int
	Comments  int
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT posts.id, posts.title, posts.content, users.pseudo, posts.created_at,
			(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id) as like_count,
			(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) as comment_count
		FROM posts
		JOIN users ON posts.user_id = users.id
		ORDER BY posts.created_at DESC`)
	if err != nil {
		log.Println("Erreur lors de la récupération des posts:", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []PostData
	for rows.Next() {
		var p PostData
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt, &p.Likes, &p.Comments)
		if err == nil {
			posts = append(posts, p)
		}
	}

	tmpl := template.Must(template.ParseFiles("_templates_/posts.html"))
	tmpl.Execute(w, posts)
}
